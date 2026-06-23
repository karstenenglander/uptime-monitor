-- name: ListSites :many
SELECT
  *
FROM
  sites;

-- name: FindSitesByID :one
SELECT * FROM sites WHERE id = $1;

-- name: RemoveSiteByID :one
DELETE FROM sites WHERE id = $1 returning name;

-- name: AddSite :one
INSERT INTO sites (name, url) VALUES ($1, $2) returning id;

-- name: UpdateSitePolled :one
UPDATE sites SET polled_at = $1 WHERE id = $2 returning id;
