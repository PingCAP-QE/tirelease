package repository

import (
	"testing"

	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestReleaseVersion(t *testing.T) {
	// Init
	var config = generateConfig()
	database.Connect(config)

	// Select
	var option = &entity.ReleaseVersionOption{}
	versions, err := SelectReleaseVersion(option)
	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(*versions) > 0)
}

/**
sql:

INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-04-02 00:00:00','4.0.12', 4, 0, 12, '', '尹苏', 'Patch', 'released', "release-4.0");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-05-28 00:00:00','4.0.13', 4, 0, 13, '', '尹苏', 'Patch', 'released', "release-4.0");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-07-27 00:00:00','4.0.14', 4, 0, 14, '', '尹苏', 'Patch', 'released', "release-4.0");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-09-27 00:00:00','4.0.15', 4, 0, 15, '', '尹苏', 'Patch', 'released', "release-4.0");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-12-17 00:00:00','4.0.16', 4, 0, 16, '', '尹苏', 'Patch', 'released', "release-4.0");

INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-04-07 00:00:00','5.0.0', 5, 0, 0, '', '尹苏', 'Major', 'released', "release-5.0");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-04-25 00:00:00','5.0.1', 5, 0, 1, '', '尹苏', 'Patch', 'released', "release-5.0");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-06-10 00:00:00','5.0.2', 5, 0, 2, '', '尹苏', 'Patch', 'released', "release-5.0");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-07-02 00:00:00','5.0.3', 5, 0, 3, '', '尹苏', 'Patch', 'released', "release-5.0");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-09-27 00:00:00','5.0.4', 5, 0, 4, '', '尹苏', 'Patch', 'released', "release-5.0");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-12-03 00:00:00','5.0.5', 5, 0, 5, '', '尹苏', 'Patch', 'released', "release-5.0");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-12-31 00:00:00','5.0.6', 5, 0, 6, '', '尹苏', 'Patch', 'released', "release-5.0");
INSERT INTO release_version (create_time, update_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '5.0.7', 5, 0, 7, '', '尹苏', 'Patch', 'upcoming', "release-5.0");

INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-06-24 00:00:00','5.1.0', 5, 1, 0, '', '尹苏', 'Minor', 'released', "release-5.1");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-07-30 00:00:00','5.1.1', 5, 1, 1, '', '尹苏', 'Patch', 'released', "release-5.1");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-09-27 00:00:00','5.1.2', 5, 1, 2, '', '尹苏', 'Patch', 'released', "release-5.1");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-12-03 00:00:00','5.1.3', 5, 1, 3, '', '尹苏', 'Patch', 'released', "release-5.1");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2022-02-22 00:00:00','5.1.4', 5, 1, 4, '', '尹苏', 'Patch', 'released', "release-5.1");
INSERT INTO release_version (create_time, update_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '5.1.5', 5, 1, 5, '', '尹苏', 'Patch', 'upcoming', "release-5.1");

INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-08-27 00:00:00','5.2.0', 5, 2, 0, '', '尹苏', 'Minor', 'released', "release-5.2");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-09-09 00:00:00','5.2.1', 5, 2, 1, '', '尹苏', 'Patch', 'released', "release-5.2");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-10-29 00:00:00','5.2.2', 5, 2, 2, '', '尹苏', 'Patch', 'released', "release-5.2");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-12-03 00:00:00','5.2.3', 5, 2, 3, '', '尹苏', 'Patch', 'released', "release-5.2");
INSERT INTO release_version (create_time, update_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '5.2.4', 5, 2, 4, '', '尹苏', 'Patch', 'upcoming', "release-5.2");

INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2021-11-30 00:00:00','5.3.0', 5, 3, 0, '', '尹苏', 'Minor', 'released', "release-5.3");
INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2022-03-03 00:00:00','5.3.1', 5, 3, 1, '', '尹苏', 'Patch', 'released', "release-5.3");

INSERT INTO release_version (create_time, update_time, actual_release_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '2022-02-15 00:00:00','5.4.0', 5, 4, 0, '', '尹苏', 'Minor', 'released', "release-5.4");
INSERT INTO release_version (create_time, update_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '5.4.1', 5, 4, 1, '', '尹苏', 'Patch', 'upcoming', "release-5.4");

INSERT INTO release_version (create_time, update_time, name, major, minor, patch, addition, owner, type, status, release_branch) VALUES (Now(), Now(), '6.0.0', 6, 0, 0, '', '尹苏', 'Major', 'upcoming', "release-6.0");

**/
