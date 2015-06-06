// Tote is a CLI application for generating structs which store SQL queries
// as defined by the directory and file structure supplied (default is
// "./sqltote").  Only .sql files are read.
//
//  Available flags:
//
//  --in={dir}          Set the SQL storage directory.  // "./sqltote"
//  --out={dir}         Set the tote package directory. // "./"
//  --file={filename}   Set the tote file name.         // "sqltote.go"
//  --pkg={package}     Set the tote package name.
//  --prefix={name}     Set the tote struct prefix.
//
// Normally, this command should be called using go:generate.  If called from
// the command line, "pkg" must be set.  The following usage will produce a
// package named "totepkg" within the "totepkg" directory:
// 	//go:generate tote -in=resources/sql/tote -out=totepkg
//
// The following usage will add a second file to the "totepkg" package:
// 	//go:generate tote -in=other/sql/tote -out=totepkg -prefix=other -file=other.go
//
// Queries are accessible in this way:
// 	import "gitreposrus.com/myproject/totepkg"
//
// 	func main() {
// 		// File originally located at "./resources/sql/tote/user/all.sql"
// 		fmt.Println(totepkg.User.All)
//
// 		// File originally located at "./resources/sql/tote/user/role/many_by_user.sql"
// 		fmt.Println(totepkg.UserRole.ManyByUser)
//
// 		// File originally located at "./other/sql/tote/user/one_by_name.sql"
// 		fmt.Println(totepkg.OtherUser.OneByName)
// 	}
//
// The main caveat seems to be naming collisions which was the primary
// motivation for the prefix flag.  Stay aware and problems can be avoided.
//
// This package started as a fork of smotes/purse.
package main
