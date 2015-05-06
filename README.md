# purse

Purse is a CLI application for generating structs which store SQL queries as
defined by the directory and file structure supplied (default is "./sqlpurse").
Only .sql files are read.

    Available flags:
    --in={dir}          Set the SQL storage directory.
    --out={dir}         Set the purse package directory.
    --file={filename}   Set the purse file name.
    --pkg={package}     Set the purse package name.
    --prefix={name}     Set the purse struct prefix.

Normally, this command should be called using go:generate. The following usage
will produce a package named "pursepkg" within the "pursepkg" directory:

    //go:generate purse -in=resources/sql/purse -out=pursepkg

The following usage will add a second file to the "pursepkg" package:

    //go:generate purse -in=other/sql/purse -out=pursepkg -prefix=other -file=other.go

Queries are accessible in this way:

    import "vcs-storage.nil/mycurrentproject/pursepkg"

    func main() {
    	// File originally located at "./resources/sql/purse/user/all.sql"
    	fmt.Println(pursepkg.User.All)

    	// File originally located at "./resources/sql/purse/user/role/many_by_user.sql"
    	fmt.Println(pursepkg.UserRole.ManyByUser)

    	// File originally located at "./other/sql/purse/user/one_by_name.sql"
    	fmt.Println(pursepkg.OtherUser.OneByName)
    }

The main caveat seems to be naming collisions which was the primary motivation
for the prefix flag. Stay aware and problems can be avoided.
