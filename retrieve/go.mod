module app/retrieve

go 1.17

replace app/bm25 => ../bm25

require (
	app/bm25 v0.0.0-00010101000000-000000000000
	app/ext v0.0.0-00010101000000-000000000000
	app/getDoc v0.0.0-00010101000000-000000000000
)

replace app/getDoc => ../getDoc

replace app/ext => ../ext
