-- Migration: Add dual language support to dokumen_kebijakan
-- Remove slug (not needed for fixed 7 pages)

-- Add dual language columns
ALTER TABLE dokumen_kebijakan
ADD COLUMN judul_en VARCHAR(100),
ADD COLUMN konten_en TEXT;

-- Update existing data with actual English content
UPDATE dokumen_kebijakan SET judul_en = 'About Us', konten_en = $$<h2>Who We Are</h2>
<p>Bulky.id is Indonesia's leading B2B wholesale re-commerce platform. We source surplus, overstock, and returned inventory directly from major e-commerce platforms, logistics providers, and corporations, process every item through our professional warehouse facility, and make it available to resellers, traders, and SMEs at transparent wholesale prices.</p>
<p>Founded in 2022 and headquartered in Jakarta, Bulky.id has processed over 700 metric tonnes of goods and serves thousands of active resellers across Indonesia. Our technology-first approach — featuring a proprietary Warehouse Management System (WMS), dedicated mobile app, and real-time inventory tracking — ensures transparency, reliability, and efficiency at every step.</p>
<h2>Vision</h2>
<p>To be Indonesia's most trusted and transparent wholesale re-commerce platform — empowering the next generation of Indonesian resellers and SMEs through reliable access to quality surplus inventory at honest, competitive prices.</p>
<h2>Mission</h2>
<p>Bulky.id exists to serve resellers, traders, and small businesses with honesty and reliability. Our mission is carried out through four commitments:</p>
<ol>
<li>To provide resellers and traders with reliable, consistent access to quality wholesale surplus inventory through a transparent, technology-enabled platform.</li>
<li>To support Indonesia's MSME ecosystem by lowering barriers to wholesale sourcing — giving anyone the real opportunity to build a viable resale business.</li>
<li>To promote responsible commercial practices by giving surplus, overstock, and returned goods a second productive life — reducing commercial and logistical waste.</li>
<li>To continuously improve our platform, warehouse operations, and service quality to deliver the most professional, honest, and efficient wholesale buying experience in Indonesia.</li>
</ol>
<h2>Our Values</h2>
<ul>
<li><strong>Transparency</strong> | Clear pricing, honest condition grades, no hidden fees. What you see is what you get.</li>
<li><strong>Reliability</strong> | Consistent inventory supply, professional warehouse operations, and committed timelines — every order, every time.</li>
<li><strong>Access</strong> | Leveling the playing field for Indonesian resellers and SMEs. Quality wholesale inventory should not be the privilege of large buyers only.</li>
<li><strong>Sustainability</strong> | Every purchase gives surplus inventory a second life. Bulky.id actively reduces commercial waste through responsible re-commerce.</li>
<li><strong>Integrity</strong> | Honest business practices, clear terms, and direct communication — because your business depends on a partner you can trust.</li>
</ul>$$ WHERE urutan = 1;
UPDATE dokumen_kebijakan SET judul_en = 'How to Buy', konten_en = $$<p>Buying at Bulky.id is easy. Follow the steps below to place your first order. For enquiries, contact our team via WhatsApp or email.</p>
<h2>Main Purchase Flow</h2>
<h3>STEP 1 — Create Your Account</h3>
<p>Download the Bulky.id app or visit bulky.id. Click Register and complete your profile. Account activation is typically within 1 x 24 hours.</p>
<h3>STEP 2 — Browse Products</h3>
<p>Once logged in, browse the catalogue by category or use the search function. Each listing shows the condition grade, quantity, and price.</p>
<h3>STEP 3 — Add to Cart</h3>
<p>Select the product and quantity, then add to cart.</p>
<h3>STEP 4 — Choose Delivery or Pickup</h3>
<p>Select self-pickup at Cibinong Warehouse, or delivery via Deliveree or Forwarder.ai. Delivery costs are calculated at checkout.</p>
<h3>STEP 5 — Choose Payment Method</h3>
<p>Select Pay Directly or Pay Patungan with Friends to share the cost with partners.</p>
<h3>STEP 6 — Complete Payment</h3>
<p>Choose your payment channel and complete within the specified time limit. Unpaid orders will be automatically cancelled.</p>
<h2>Self-Pickup Information</h2>
<p>Warehouse: Jl. Raya Mayor Oking No. 62a, Cibinong, West Java. Monday to Saturday, 09:00 to 18:00 WIB.</p>$$ WHERE urutan = 2;
UPDATE dokumen_kebijakan SET judul_en = 'About Payment', konten_en = $$<h2>Virtual Accounts (Bank Transfer)</h2>
<p>Supported banks: BCA, BNI, BRI, Mandiri, Permata, BSI, CIMB Niaga, Danamon, BJB.</p>
<h2>Credit &amp; Debit Cards</h2>
<p>Visa, Mastercard, JCB.</p>
<h2>E-Wallets</h2>
<p>OVO, DANA, GoPay, ShopeePay, LinkAja.</p>
<h2>QRIS</h2>
<p>Scan the QRIS code for instant payment from any mobile banking app or e-wallet that supports QRIS.</p>
<p>Note: No cash or COD. All payments are processed via Xendit.</p>$$ WHERE urutan = 3;
UPDATE dokumen_kebijakan SET judul_en = 'Contact Us', konten_en = $$<p>We are here to help. Contact us through:</p>
<h2>Warehouse Address</h2>
<p>Jl. Raya Mayor Oking Jaya Atmaja No.62a, Kel Cirimekar, Kec. Cibinong, Kabupaten Bogor, West Java 16918</p>
<h2>Phone / WhatsApp</h2>
<p><a href="https://wa.me/62811833164">+62811 833 164</a></p>
<h2>Email</h2>
<p><a href="mailto:admin@bulky.id">admin@bulky.id</a></p>
<h2>Operating Hours</h2>
<p>Monday to Saturday: 09:00 to 18:00 WIB. Sunday: Closed.</p>$$ WHERE urutan = 4;
UPDATE dokumen_kebijakan SET judul_en = 'Frequently Asked Questions', konten_en = $$<p>Can't find what you're looking for? Contact us via WhatsApp at <a href="https://wa.me/62811833164">0811 833 164</a> or email <a href="mailto:admin@bulky.id">admin@bulky.id</a>.</p>
<h2>About Bulky.id</h2>
<h3>What is Bulky.id?</h3>
<p>Bulky.id is Indonesia's B2B wholesale re-commerce platform. We source surplus, overstock, and returned inventory from e-commerce platforms, logistics providers, and corporations, then make it available to resellers, traders, and SMEs at wholesale prices.</p>
<h3>Is there a minimum order?</h3>
<p>There is no fixed minimum order value. Products are sold in pre-packed bulk units. The minimum purchase is the smallest available unit for each listing.</p>
<h2>About Products</h2>
<h3>What are the condition grades?</h3>
<p>New (factory-sealed), Like New (opened but unused), Used (previously used, functional), and Salvage (damaged, for parts or recycling).</p>
<h2>About Ordering and Payment</h2>
<h3>What payment methods are accepted?</h3>
<p>Virtual Account, Credit/Debit Cards, E-Wallets (OVO, DANA, GoPay, ShopeePay, LinkAja), and QRIS via Xendit.</p>
<h3>What is Patungan?</h3>
<p>Bulky.id's co-purchase split-payment feature. Two or more registered buyers split the cost of one order.</p>$$ WHERE urutan = 5;
UPDATE dokumen_kebijakan SET judul_en = 'Terms and Conditions', konten_en = $$<p>The full Buyer Terms and Conditions govern all transactions on Bulky.id. By purchasing on Bulky.id, you agree to be bound by the full Terms and Conditions.</p>
<h2>Key Principles</h2>
<ul>
<li><strong>Platform Purpose:</strong> B2B wholesale re-commerce for resellers, traders, and SMEs — not a consumer marketplace.</li>
<li><strong>Product Nature:</strong> All products are surplus, overstock, or returned goods sold on an AS-IS basis.</li>
<li><strong>No Warranties:</strong> Bulky.id disclaims all implied warranties.</li>
<li><strong>All Sales Final:</strong> Returns not permitted except for material quantity shortfalls reported within 24 hours.</li>
<li><strong>Governing Law:</strong> Indonesian law. Disputes via South Jakarta District Court.</li>
</ul>$$ WHERE urutan = 6;
UPDATE dokumen_kebijakan SET judul_en = 'Privacy Policy', konten_en = $$<p>Bulky.id is committed to protecting the privacy and personal data of all users. This Privacy Policy describes how we collect, use, store, share, and protect your personal information.</p>
<h2>Personal Data We Collect</h2>
<ul>
<li>Account registration: full name, email, phone, business information</li>
<li>Transaction data: order history, payment records, shipping addresses</li>
<li>Usage data: pages visited, device and access information</li>
</ul>
<h2>Your Rights</h2>
<p>You have the right to access, correct, delete, and port your personal data. Contact <a href="mailto:admin@bulky.id">admin@bulky.id</a> with subject "Data Privacy Request".</p>$$ WHERE urutan = 7;

-- Set NOT NULL after data update
ALTER TABLE dokumen_kebijakan
ALTER COLUMN judul_en SET NOT NULL,
ALTER COLUMN konten_en SET NOT NULL;

-- Remove slug column (not needed for fixed pages)
ALTER TABLE dokumen_kebijakan
DROP COLUMN IF EXISTS slug;

-- Add comments
COMMENT ON COLUMN dokumen_kebijakan.judul IS 'Judul dalam Bahasa Indonesia';
COMMENT ON COLUMN dokumen_kebijakan.judul_en IS 'Judul dalam Bahasa Inggris';
COMMENT ON COLUMN dokumen_kebijakan.konten IS 'Konten HTML dalam Bahasa Indonesia';
COMMENT ON COLUMN dokumen_kebijakan.konten_en IS 'Konten HTML dalam Bahasa Inggris';
