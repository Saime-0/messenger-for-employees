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
		RETURNING tag_id
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
	).Err()
	return
}

func (r *TagsRepo) UpdateTag(tag *request_models.UpdateTag) (err error) {
	err = r.db.QueryRow(`
		UPDATE tags 
		SET name = $2
		WHERE tag_id = $1
	`,
		tag.TagID,
		tag.Name,
	).Err()
	return
}

func (r *TagsRepo) DropTag(tag *request_models.DropTag) (err error) {
	err = r.db.QueryRow(`
		DELETE FROM tags
		WHERE tag_id = $1
	`,
		tag.TagID,
	).Err()
	return
}

func (r *TagsRepo) TagExists(tagID int) (exists bool, err error) {
	err = r.db.QueryRow(`
		SELECT EXISTS(
		    SELECT 1
		    FROm tags
		    WHERE tag_id = $1
		)
	`,
		tagID,
	).Err()
	return
}

func (r *TagsRepo) GiveTag(tag *request_models.GiveTag) (err error) {
	err = r.db.QueryRow(`
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
	).Err()
	return
}

func (r *TagsRepo) TakeTag(tag *request_models.TakeTag) (err error) {
	err = r.db.QueryRow(`
		DELETE FROM positions
		WHERE emp_id = $2 AND tag_id = $1
	`,
		tag.TagID,
		tag.EmpID,
	).Err()
	return
}

func (r *TagsRepo) Tags(params *model.Params) (*model.Tags, error) {

	tags := &model.Tags{
		Tags: []*model.Tag{},
	}

	var rows, err = r.db.Query(`
		SELECT t.tag_id, t.name
		FROM tags t
		LIMIT $1
		OFFSET $2
		`,
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
