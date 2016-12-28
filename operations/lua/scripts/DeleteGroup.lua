local userHashKey = KEYS[1]
local groupIDKey = KEYS[2]
local groupHashKey = KEYS[3]
local groupMembersKey = KEYS[4]
local groupMessagesKey = KEYS[5]
local groupGeoKey = KEYS[6]
local groupLocationsKey = KEYS[7]

local userID = ARGV[1]
local groupID = ARGV[2]

local userGroupsKeyPartial = ARGV[3]
local userGrpMsgKeyPartial = ARGV[4]

-- check if user is group owner
if redis.call("HGET", groupHashKey, "owner") ~= userID then
    return redis.error_reply("Invalid permissions.")
end

-- get list of ID's of group members
local groupMemberIDList = redis.call("HKEYS", groupMembersKey)

-- delete group from each user's group list
for i = 1, #groupMemberIDList do
	-- delete group from user group hash set
	redis.call("HDEL", userGroupsKeyPartial .. groupMemberIDList[i], groupID)
	-- delete messages from user group messages
	-- maybe change this in the future because the message get
	-- deleted after ttl is up anyway
	redis.call("DEL", userGrpMsgKeyPartial .. groupMemberIDList[i] .. ":" .. groupID)
end

local groupName = redis.call("HGET", groupHashKey, "groupName")

redis.call("DEL", groupHashKey)
redis.call("DEL", groupMembersKey)
redis.call("DEL", groupMessagesKey)
redis.call("DEL", groupGeoKey)
redis.call("DEL", groupLocationsKey)
redis.call("HINCRBY", userHashKey, "adminGroupCount", -1)
redis.call("HDEL", groupIDKey, groupName)

return "Success"
