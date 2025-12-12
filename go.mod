module github.com/absfs/rofs

go 1.23

require (
	github.com/absfs/absfs v0.0.0-20251208232938-aa0ca30de832
	github.com/absfs/ioutil v0.0.0-20180730202941-e848c1afe5fd
	github.com/absfs/memfs v0.0.0-20251122223403-0826c6ce3884
	github.com/absfs/osfs v0.1.0-fastwalk
)

require github.com/absfs/inode v0.0.0-20251208170702-9db24ab95ae4 // indirect

replace (
	github.com/absfs/absfs => ../absfs
	github.com/absfs/fstools => ../fstools
	github.com/absfs/inode => ../inode
	github.com/absfs/memfs => ../memfs
	github.com/absfs/osfs => ../osfs
)
