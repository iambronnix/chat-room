simple login backend code in Go.
getLogins() captures user-detals directly from the terminal using stdin sends through 2 different channels.
User-details are stored into User{} struct.
queryUsers() scans through the available "hashes" inthe underlying database and updates the slice inthe UserHash{} struct.
Both getLogins and queryLogins are added to waitgroup to synchronize their execution.
checkUser() parses the user-details to hashFunction for hashing and then returns a hash.
a for loop ranges through the slice of hashes extracted from the database and compares each to the hashed user-details.
if a matching hash is found a true bool s returned.else false is returned.
in finalFlow() if true is returned the user is notified else an account is created by calling insertUser().
ignore spinnerWrapper....it's just a go routine that gives the impresson that all the processes in the background are running concurrently.

