package trackernode

// SQLCreateDataNodesTable SQL to create the DataNodes table
const SQLCreateDataNodesTable string = `
	CREATE TABLE datanodes (
	id SERIAL PRIMARY KEY,
	dataNodeID int UNIQUE NOT NULL,
	ip varchar(60) NOT NULL,
	basePort varchar(60) NOT NULL
	);
`

// SQLCreateMetaFile SQL to create the Meta Files table
const SQLCreateMetaFile string = `
	CREATE TABLE metafiles (
	id SERIAL PRIMARY KEY,
	fileName varchar(60) NOT NULL,
	clientID int NOT NULL,
	fileSize int NOT NULL, 
	location varchar(60) NOT NULL
	);
	ALTER TABLE metafiles
	ADD CONSTRAINT unq_filename_clientid UNIQUE(fileName, clientID);
`

// SQLDropDataNodesTable SQL to drop the DataNodes table
const SQLDropDataNodesTable string = `DROP TABLE IF EXISTS datanodes;`

// SQLDropMetaFileTable SQL to drop the Meta Files table
const SQLDropMetaFileTable string = `DROP TABLE IF EXISTS metafiles;`

// sqlInsertDataNode SQL to insert a data node in the DataNodes table
const sqlInsertDataNode string = `
	INSERT INTO datanodes (dataNodeID, ip, basePort)
	VALUES ($1, $2, $3)
`

// sqlInsertFileEntry SQL to insert a a file entry into the Meta File table
const sqlInsertFileEntry string = `
	INSERT INTO metafiles (fileName, clientID, fileSize, location)
	VALUES ($1, $2, $3, $4)
`

// sqlDeleteDataNode SQL to delete a data node from the DataNodes table
const sqlDeleteDataNode string = `DELETE FROM datanodes WHERE dataNodeID=$1`

// sqlSelectAllDataNodes SQL to select all datanodes
const sqlSelectAllDataNodes string = `SELECT * FROM datanodes`

// sqlSelectDataNode SQL to select a datanode indentified by its ID
const sqlSelectDataNode string = `
	SELECT * FROM datanodes WHERE dataNodeID=$1
`

// sqlSelectAllMetaFiles SQL to select all entries in the metafile table
const sqlSelectAllMetaFiles string = `SELECT * FROM metafiles`

// sqlSelectMetaFile SQL to select a meta file entry
const sqlSelectMetaFile string = `
	SELECT * FROM metafiles WHERE fileName = $1 and clientID = $2
`

// sqlUpdateMetaFile SQL to update a metafile entry
const sqlUpdateMetaFile string = `
	UPDATE metafiles
	SET location = $1
	WHERE fileName = $2 AND clientID = $3; 
`

// sqlSelectAllMetaFilesForClient SQL to select all metafiles for a client
const sqlSelectAllMetaFilesForClient string = `
	SELECT fileName, fileSize 
	FROM metafiles 
	WHERE clientID = $1
`
