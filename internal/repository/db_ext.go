package repository

func (q *Queries) DB() DBTX {
	return q.db
}
