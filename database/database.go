package database

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"

	. "../network/structures"
	_ "github.com/go-sql-driver/mysql"
)

// Database wrapper class

var DB *sql.DB = nil

// Open "opens" a database connection
// It returns the current instance of the DB object
func Open() *sql.DB {
	if DB != nil {
		return DB
	}

	var err error
	DB, err = sql.Open("mysql",
		"root:Nz2y@-NycwT%2e4s@tcp/soulworker")

	if err != nil {
		log.Fatal(err)
	}

	return DB
}

// CanConnect checks whether or not a database is accepting connections
// It returns an error if not
func CanConnect() error {
	fmt.Print("[Database]\tChecking connection... ")
	if err := DB.Ping(); err != nil {
		fmt.Println("Could not connect: ", err)
		return err
	} else {
		fmt.Println("Connected!")
	}
	return nil
}

// VerifyLoginCredentials is used to determine whether a given username and password combination correlates to a given account
// It returns the account id if valid and 0 if invalid
func VerifyLoginCredentials(username string, password string) uint32 {

	var (
		id         int
		dbUsername string
		dbPassword string
	)

	err := DB.QueryRow("SELECT id, username, password from users where username = ?", username).Scan(&id, &dbUsername, &dbPassword)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	checksum := sha256.Sum256([]byte(password))
	comparison := fmt.Sprintf("%x", checksum)

	fmt.Printf("[+] VerifyLoginCredentials:\n\tId: %d\n\tUsername: %s\n\tPassword: %s\n\tPassword hash: %s\n\tMatch: %v\n\n", id, dbUsername, dbPassword, comparison, dbPassword == comparison)

	if comparison != dbPassword {
		return 0
	}

	return uint32(id)
}

// Updates a given accountId's sessionKey in the database
func UpdateSessionKey(accountId uint32, sessionKey uint64) {
	stmt, err := DB.Prepare("UPDATE users SET session_key = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(int64(sessionKey), accountId)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
}

// Verifies that an accountId's sessionKey matches the one in the database
func VerifySessionKey(accountId uint32, sessionKey uint64) int {

	var id int

	err := DB.QueryRow("SELECT id from users where id = ? and session_key = ?", int(accountId), int64(sessionKey)).Scan(&id)
	if err != nil {
		log.Println("Account id:", accountId, int(accountId), "| Session key:", sessionKey, int64(sessionKey))
		log.Fatal(err)
		return 0
	}

	return id
}

// Inserts an entry for the character table into the database from the contents of the CharacterInfo struct
func InsertCharacterToDb(charInfo *CharacterInfo) int64 {

	stmt, err := DB.Prepare("INSERT INTO characters(accountId, `index`, name, class, appearance, level) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Inserting character for account id:", charInfo.AccountId)
	res, err := stmt.Exec(charInfo.AccountId, charInfo.Index, charInfo.Username, charInfo.CharSelection, charInfo.Appearance, charInfo.Level)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

	return lastId
}

// Fetches the number of characters the given accountId has on the server.
func FetchUserCharacterCount(accountId uint32) int {
	var count int

	err := DB.QueryRow("SELECT COUNT(id) from characters where accountId = ?", accountId).Scan(&count)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	return count
}

// Fetches the user character with the given characterId on the given accountId.
func FetchUserCharacter(accountId uint32, characterId uint32) (CharacterInfo, error) {

	char := CharacterInfo{}

	err := DB.QueryRow("SELECT id, accountId, `index`, name, class, appearance, level from characters where id = ?", int(characterId)).Scan(&char.Id, &char.AccountId, &char.Index, &char.Username, &char.CharSelection, &char.Appearance, &char.Level)
	if err != nil {
		log.Fatal(err)
		return char, err
	}

	return char, nil
}

// Fetches the user character with the given characterId.
func FetchUserCharacterByCharacterId(characterId uint32) (CharacterInfo, error) {

	char := CharacterInfo{}

	log.Println("Character id", characterId, int(characterId))

	err := DB.QueryRow("SELECT id, accountId, `index`, name, class, appearance, level from characters where id = ?", int(characterId)).Scan(&char.Id, &char.AccountId, &char.Index, &char.Username, &char.CharSelection, &char.Appearance, &char.Level)
	if err != nil {
		log.Fatal(err)
		return char, err
	}

	return char, nil
}
