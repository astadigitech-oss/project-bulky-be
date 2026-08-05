package main

import (
	"database/sql"
	"strings"
	"time"
)

type v1User struct {
	ID             string
	Name           string
	Email          string
	Username       sql.NullString
	Phone          sql.NullString
	ProfilePicture sql.NullString
	Password       sql.NullString
	IsAdmin        bool
	DeletedAt      sql.NullTime
	CreatedAt      sql.NullTime
	UpdatedAt      sql.NullTime
}

// phaseBuyer memigrasi users -> buyer (is_admin=0) dan users(is_admin=1)+admins -> admin. Dokumen §6.
func (a *App) phaseBuyer() error {
	users, err := a.readUsers()
	if err != nil {
		return err
	}

	// Pisahkan tujuan (keputusan #11)
	var buyers, adminUsers []v1User
	for _, u := range users {
		if u.IsAdmin {
			adminUsers = append(adminUsers, u)
		} else {
			buyers = append(buyers, u)
		}
	}

	// Dedup telepon (keputusan #4): per nomor ternormalisasi, akun ter-update terbaru menang.
	phoneWinner := map[string]string{} // nomor -> user id pemenang
	winnerTime := map[string]time.Time{}
	for _, u := range buyers {
		phone, ok := normalizePhone(u.Phone.String)
		if !ok {
			continue
		}
		t := latestActivity(u)
		if cur, exists := winnerTime[phone]; !exists || t.After(cur) {
			phoneWinner[phone] = u.ID
			winnerTime[phone] = t
		}
	}

	for id := range a.tgt.BuyerIDs {
		a.buyerKnown[id] = true
	}
	batchEmails := map[string]bool{}
	batchUsernames := map[string]bool{}
	batchTelepon := map[string]bool{}

	for _, u := range buyers {
		if a.tgt.BuyerIDs[u.ID] {
			a.buyerKnown[u.ID] = true
			a.rep.Count("buyer.sudah_ada_di_v2")
			continue
		}

		// telepon
		telepon := placeholderPhone(u.ID, func(t string) bool {
			return a.tgt.BuyerTelepon[t] || batchTelepon[t]
		})
		if phone, ok := normalizePhone(u.Phone.String); ok {
			switch {
			case phoneWinner[phone] != u.ID:
				a.rep.Add("3", "buyer", u.ID, "telepon_duplikat", phone+" dipegang akun lebih baru -> placeholder")
			case a.tgt.BuyerTelepon[phone] || batchTelepon[phone]:
				a.rep.Add("3", "buyer", u.ID, "telepon_bentrok_target", phone+" sudah ada di v2 -> placeholder")
			default:
				telepon = phone
			}
		} else if strings.TrimSpace(u.Phone.String) != "" {
			a.rep.Add("3", "buyer", u.ID, "telepon_invalid", "'"+u.Phone.String+"' -> placeholder")
		} else {
			a.rep.Add("3", "buyer", u.ID, "telepon_kosong", "-> placeholder")
		}
		batchTelepon[telepon] = true

		// email (lowercase; unik di v2 termasuk baris terhapus)
		var email *string
		if e := strings.ToLower(strings.TrimSpace(u.Email)); e != "" {
			if a.tgt.BuyerEmails[e] || batchEmails[e] {
				a.rep.Add("3", "buyer", u.ID, "email_bentrok", e+" sudah dipakai -> NULL")
			} else {
				batchEmails[e] = true
				email = &e
			}
		}

		// username (v2 varchar(50) unik, nullable)
		var username *string
		if v := strings.TrimSpace(u.Username.String); v != "" {
			v = truncate(v, 50)
			if a.tgt.BuyerUsernames[v] || batchUsernames[v] {
				a.rep.Add("3", "buyer", u.ID, "username_bentrok", v+" -> NULL")
			} else {
				batchUsernames[v] = true
				username = &v
			}
		}

		nama := truncate(firstNonEmpty(u.Name, "Buyer "+id8(u.ID)), 100)
		if len([]rune(u.Name)) > 100 {
			a.rep.Add("3", "buyer", u.ID, "nama_terpotong", "")
		}

		var password *string
		if u.Password.Valid && u.Password.String != "" {
			password = strPtr(convertBcrypt(u.Password.String))
		}
		var fotoURL *string
		if f := strings.TrimSpace(u.ProfilePicture.String); f != "" {
			fotoURL = strPtr(strings.TrimPrefix(f, "public/"))
		}

		if err := a.exec(`INSERT INTO buyer (id, nama, username, email, password, telepon, foto_url,
				is_active, is_verified, telepon_verified_at, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, true, false, NULL, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			u.ID, nama, username, email, password, telepon, fotoURL,
			nullTimePtr(u.CreatedAt), nullTimePtr(u.UpdatedAt), nullTimePtr(u.DeletedAt)); err != nil {
			return err
		}
		a.buyerKnown[u.ID] = true
		a.rep.Count("buyer.insert")
	}

	return a.migrateAdmins(adminUsers)
}

// §6.1: tabel admins v1 + users(is_admin=1) -> admin v2 (role ADMIN)
func (a *App) migrateAdmins(adminUsers []v1User) error {
	batchEmails := map[string]bool{}

	insertAdmin := func(id, name, email string, password sql.NullString, createdAt, updatedAt, deletedAt sql.NullTime, source string) error {
		if a.tgt.AdminIDs[id] {
			a.rep.Count("admin.sudah_ada_di_v2")
			return nil
		}
		e := strings.ToLower(strings.TrimSpace(email))
		if e == "" {
			a.rep.Add("3", "admin", id, "skip", source+": email kosong")
			return nil
		}
		if a.tgt.AdminEmails[e] || batchEmails[e] {
			a.rep.Add("3", "admin", id, "email_bentrok", source+": "+e+" sudah dipakai -> skip")
			return nil
		}
		pwd := "!migrated-no-password"
		if password.Valid && password.String != "" {
			pwd = convertBcrypt(password.String)
		} else {
			a.rep.Add("3", "admin", id, "tanpa_password", source+": password kosong -> login dinonaktifkan, reset manual")
		}
		if err := a.exec(`INSERT INTO admin (id, nama, email, password, role_id, is_active, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, true, ?, ?, ?) ON CONFLICT (id) DO NOTHING`,
			id, truncate(firstNonEmpty(name, "Admin "+id8(id)), 100), e, pwd, a.tgt.RoleAdminID,
			nullTimePtr(createdAt), nullTimePtr(updatedAt), nullTimePtr(deletedAt)); err != nil {
			return err
		}
		batchEmails[e] = true
		a.rep.Count("admin.insert")
		return nil
	}

	// tabel admins v1 duluan (admin "asli"), lalu users is_admin=1
	rows, err := a.my.Query(`SELECT id, name, email, password, deleted_at, created_at, updated_at FROM admins ORDER BY created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id, name, email string
		var password sql.NullString
		var deletedAt, createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&id, &name, &email, &password, &deletedAt, &createdAt, &updatedAt); err != nil {
			return err
		}
		if err := insertAdmin(id, name, email, password, createdAt, updatedAt, deletedAt, "admins"); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, u := range adminUsers {
		if err := insertAdmin(u.ID, u.Name, u.Email, u.Password, u.CreatedAt, u.UpdatedAt, u.DeletedAt, "users.is_admin"); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) readUsers() ([]v1User, error) {
	rows, err := a.my.Query(`SELECT id, name, email, username, phone_number, profile_picture,
			password, is_admin, deleted_at, created_at, updated_at
		FROM users ORDER BY created_at, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []v1User
	for rows.Next() {
		var u v1User
		var isAdmin sql.NullInt64
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Username, &u.Phone, &u.ProfilePicture,
			&u.Password, &isAdmin, &u.DeletedAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		u.IsAdmin = boolFromTiny(isAdmin)
		out = append(out, u)
	}
	return out, rows.Err()
}

func latestActivity(u v1User) time.Time {
	if u.UpdatedAt.Valid {
		return u.UpdatedAt.Time
	}
	if u.CreatedAt.Valid {
		return u.CreatedAt.Time
	}
	return time.Time{}
}
