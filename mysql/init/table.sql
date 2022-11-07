
/* CREATE USER 'pendex'@'172.19.0.3' identified by 'dedicatus545';
CREATE DATABASE JohannesIndex; 
*/
USE JohannesIndex;
CREATE TABLE m_user (
    UserID varchar(50) not null primary key,
    Password varchar(1000) not null,
    DisplayName varchar(20) not null,
    AccountLevel int,
    CreateDate date,
    UpdateDate date,
    ESIndexName varchar(30) not null unique,
    index M_User_Index (UserID)
);

INSERT INTO m_user(UserID, Password, DisplayName, AccountLevel,  ESIndexName) VALUES 
("admin","ede3518e854d8dc30965be15005f4ae7845a610bf171145a4e2473752d654740", "admin", "3",  "bookmark-uuid16");

CREATE TABLE m_folder (
    FolderID varchar(50) not null primary key,
    UserID varchar(50) not null,
    FolderName varchar(20) not null, 
    index U_Folder_Index (FolderID),
    foreign key fku_M_User_ID(UserID) references m_user(UserID)
);
/*ALTER TABLE M_Folder MODIFY COLUMN ParentID varchar(20);*/
INSERT INTO m_folder (FolderID, UserID, FolderName)
VALUES ("63b96ec6-e9d7-4fa4", "admin", "全体");
INSERT INTO m_folder (FolderID, UserID, FolderName)
VALUES ("8160d492-1ead-4d33", "admin", "github");

/* grant create on JohannesIndex.M_Folder to 'pendex'@'172.19.0.3';
grant create on JohannesIndex.M_User to 'pendex'@'172.19.0.3';
grant select, insert, update, delete on JohannesIndex.M_Folder to 'pendex'@'172.19.0.3';
grant select, insert, update, delete on JohannesIndex.M_User to 'pendex'@'172.19.0.3';
*/