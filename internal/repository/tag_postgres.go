package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/admin/request_models"
)

type TagsRepo struct {
	db *sql.DB
}

func NewTagsRepo(db *sql.DB) *TagsRepo {
	return &TagsRepo{
		db: db,
	}
}

func (r *TagsRepo) CreateTag(tag *request_models.CreateTag) (tagID int, err error) {
	err = r.db.QueryRow(`
		INSERT INTO tags (name) VALUES ($1)
		RETURNING id
	`,
		tag.Name,
	).Scan(&tagID)
	return
}

func (r *TagsRepo) TagExistsByName(name string) (exists bool, err error) {
	err = r.db.QueryRow(`
		SELECT EXISTS(
		    SELECT 1
		    FROm tags
		    WHERE name = $1
		)
	`,
		name,
	).Scan(&exists)
	return
}

func (r *TagsRepo) UpdateTag(tag *request_models.UpdateTag) (err error) {
	_, err = r.db.Exec(`
		UPDATE tags 
		SET name = $2
		WHERE id = $1
	`,
		tag.TagID,
		tag.Name,
	)
	return
}

func (r *TagsRepo) DropTag(tag *request_models.DropTag) (err error) {
	_, err = r.db.Exec(`
		DELETE FROM tags
		WHERE id = $1
	`,
		tag.TagID,
	)
	return
}

func (r *TagsRepo) TagExistsByID(tagID int) (exists bool, err error) {
	err = r.db.QueryRow(`
		SELECT EXISTS(
		    SELECT 1
		    FROm tags
		    WHERE id = $1
		)
	`,
		tagID,
	).Scan(&exists)
	return
}

func (r *TagsRepo) GiveTag(tag *request_models.GiveTag) (err error) {
	_, err = r.db.Exec(`
		WITH "except"(tag_id) AS (
		    SELECT tag_id
		    FROM positions
		    WHERE tag_id = ANY($1) AND emp_id = $2
		)
		INSERT INTO positions (tag_id, emp_id) 
		SELECT tagid, $2
		FROM unnest($1::bigint[]) inp(tagid)
		WHERE tagid != ALL(select tag_id from "except")
	`,
		pq.Array(tag.TagIDs),
		tag.EmpID,
	)
	return
}

func (r *TagsRepo) TakeTag(tag *request_models.TakeTag) (err error) {
	_, err = r.db.Exec(`
		DELETE FROM positions
		WHERE emp_id = $2 AND tag_id = $1
	`,
		tag.TagID,
		tag.EmpID,
	)
	return
}

func (r *TagsRepo) Tags(tagIDs []int, params *model.Params) (*model.Tags, error) {
	var (
		allRows bool
		tags    = &model.Tags{
			Tags: []*model.Tag{},
		}
	)
	if len(tagIDs) == 0 {
		allRows = true
	}
	var rows, err = r.db.Query(`
		SELECT id, name
		FROM tags
		WHERE (
		    $1::BOOLEAN
		    OR
		    id = ANY($2::BIGINT[])
		)
		LIMIT $3
		OFFSET $4
		`,
		allRows,
		pq.Array(tagIDs),
		params.Limit,
		params.Offset,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := new(model.Tag)
		if err = rows.Scan(&m.TagID, &m.Name); err != nil {
			return nil, err
		}
		tags.Tags = append(tags.Tags, m)
	}

	return tags, nil
}
