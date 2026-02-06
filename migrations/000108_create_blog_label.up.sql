CREATE TABLE blog_label (
    blog_id UUID NOT NULL,
    label_id UUID NOT NULL,
    CONSTRAINT pk_blog_label PRIMARY KEY (blog_id, label_id),
    CONSTRAINT fk_blog_label_blog FOREIGN KEY (blog_id) REFERENCES blog(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog_label_label FOREIGN KEY (label_id) REFERENCES label_blog(id) ON DELETE CASCADE
);

CREATE INDEX idx_blog_label_blog ON blog_label(blog_id);
CREATE INDEX idx_blog_label_label ON blog_label(label_id);
