package oxr

import "fmt"

var (
    VersionMajor = 0
    VersionMinor = 1
    VersionRevision = 0
    VersionTag = ""
    Version = fmt.Sprintf("%v.%v.%v", VersionMajor, VersionMinor, VersionRevision)
    AuthorName = "Martín Raúl Villalba"
    AuthorEMail = "martin@martinvillalba.com"
    Author = fmt.Sprintf("%v <%v>", AuthorName, AuthorEMail)
)

func init() {
    if VersionTag != "" {
        Version = fmt.Sprintf("%v-%v", Version, VersionTag)
    }
}
