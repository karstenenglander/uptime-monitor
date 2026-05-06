-- name: ListSites :many
SELECT
  *
FROM
  sites;

-- name: FindSitesByID :one
SELECT * FROM sites WHERE id = $1;