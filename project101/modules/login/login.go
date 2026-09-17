package login

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type User struct {
	username string
	password string
}
type UserHash struct {
	getUserHash []string
}

var (
	loginName     = make(chan string, 1)
	loginPassword = make(chan string, 1)
	userName = make(chan string, 1)
)



func getLogins(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("\nEnter username:")
	terminalInput := os.Stdin
	scanner := bufio.NewScanner(terminalInput)
	for scanner.Scan() {
		loginName <- scanner.Text()
		close(loginName)
		break
	}
	fmt.Println("Enter password:")
	for scanner.Scan() {
		loginPassword <- scanner.Text()
		close(loginPassword)
		break
	}
	if userErr := scanner.Err(); userErr != nil {
		log.Println("Error capturing user details") //recover if users logins fail to be inserted
	}

}
//this func queries users crypted hash from the database and returns a slice of users hashes and error which is either nil 
func queryUsers(wg *sync.WaitGroup, db *sql.DB)( []string, error) {
	
	defer wg.Done()
	var hash string
	h := &UserHash{}
	scanner, prepErr := db.Prepare("SELECT userHash FROM users") //prepare a query mainly to prevent sql injection
	if prepErr != nil {
		return nil, prepErr
	}
	results, queryErr := scanner.Query()
	if queryErr != nil {
		return nil, queryErr
	}

	for results.Next() { //scan through the available hashes

		scanErr := results.Scan(&hash)
		if scanErr != nil {
			return nil, scanErr
		}
		h.getUserHash = append(h.getUserHash, hash) //continously update the UserHash{} struct

		//fmt.Println(hash) //send the queries through a channel....not a good idea though
	}
	return h.getUserHash, nil

}
func hashFunction(text string) string { //converts a parsed text into a "hash"
	//n := newUser()
	//userString := n.username + n.password
	resultHash := make([]byte, len(text))
	for i := 0; i < len(text); i++ {
		char := text[i]
		switch {
		case char >= 'a' && char <= 'z':
			resultHash[i] = 'a' + (char-'a'+13)%26
		case char >= 'A' && char <= 'Z':
			resultHash[i] = 'A' + (char-'A'+13)%26
		default:
			resultHash[i] = char
		}

	}
	return string(resultHash) //combines back the characters into strings
	//this is basically rot13
	// still figuring how deal with numbers
}

func (s *User) checkUser(dbHash []string) bool { //*userHash []string**value to pass into the function

	stringHash := s.username + s.password + "iamb37!!45ronnix"
	userRot := hashFunction(stringHash)
	for _, i := range dbHash {
		if i == userRot{
			return true
		}else{
			continue
		}
	}
	return false
}
func (n *User) insertUser(db *sql.DB) {
	newUserHash := hashFunction(n.username + n.password + "iamb37!!45ronnix")               //create a hash for the new user
	newQuery, prepErr := db.Prepare("INSERT INTO users VALUES($1,$2)") //restrict variables to be inserted
	if prepErr != nil {
		panic(prepErr)
	}

	_, queryErr := newQuery.Exec(newUserHash, n.username) //inserts a new user to the table
	if queryErr != nil {
		panic(queryErr)
	}
	defer db.Close()

}
//this func drops the whole table and returns error otherwise
func dropTable(db *sql.DB) error{ 
	_, dropErr := db.Exec("DROP TABLE users")
	if dropErr != nil {
		return dropErr
	}
	return nil
}
func createTable(db *sql.DB){ //use to create a new table
	//if you create a new table remember to change table name
	defer func(){
		if recovErr := recover();recovErr!=nil{
			//intentionally left it blank
		}
	}()
	createStatement := `
	CREATE TABLE users(
	userHash text NOT NULL UNIQUE,
	userName text
	)
	WITH(
	OIDS=FALSE
	)
	TABLESPACE pg_default;
	ALTER TABLE users
	OWNER to postgres;
	`
	_, createErr := db.Exec(createStatement)
	if createErr != nil {
		panic(createErr)
	}
}

func dupChecker(hash []string) {
	counts := make(map[string]int)
	for _, hashString := range hash {
		writeErr := os.WriteFile("hashFile.txt", []byte(hashString), 0644)
		if writeErr != nil {
			fmt.Fprintf(os.Stderr, "\t%v\n", writeErr)
		}
	}
	hashFile, openErr := os.Open("hashFile.txt")
	if openErr != nil {
		fmt.Fprintf(os.Stderr, "\t%v\n", openErr)

	}
	input := bufio.NewScanner(hashFile)
	for input.Scan() {
		if eofErr := input.Err(); eofErr == io.EOF {
			break
		}
		counts[input.Text()]++
	}
	for line, n := range counts {
		if 1 <= n {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}

}
func (newUser *User) finalFlow(dbHashes []string, db *sql.DB){
	userExists := newUser.checkUser(dbHashes)
	switch userExists{
		case true:
		userName <- newUser.username
		fmt.Println("**********Account Exists************\nusername:", newUser.username)
		case false:
		fmt.Println("***********Creating Account************")
		time.Sleep(5 * time.Second)
		fmt.Println("Done!!!")
		newUser.insertUser(db)
		userName <- newUser.username
		default:
		log.Fatal("error getting user details")
	}

}

func Login(db *sql.DB)(error, string){
	//dropTable()//avoid as much as you can
	createTable(db)
	defer func() { //recover panics
		if recoverErr := recover(); recoverErr != nil {
			fmt.Fprintf(os.Stderr, "%v\t\n", recoverErr)
		}

	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	getLogins(wg)
	newUser := User{
		username: <-loginName,
		password: <-loginPassword,
	}
  	dbHashes, queryErr := queryUsers(wg,db)
      if queryErr != nil{
        return queryErr, fmt.Sprintf("%v",queryErr)     
      }
     wg.Wait()
	
	newUser.finalFlow(dbHashes,db)
	// dupChecker(dbHashes)//there's a logical bug in this function
	return  nil, <-userName

}
