local userHashKey = KEYS[1]
local userIDKey = KEYS[2]
local userGroupsKey = KEYS[3]
local userGroupInvitesKey = KEYS[4]

local userID = ARGV[1]

local userGrpMsgKeyPartial = ARGV[2]
local groupMemKeyPartial = ARGV[3]
local groupLocationsKeyPartial = ARGV[4]

-- check if user exists
if redis.call("EXISTS", userHashKey) == 0 then
	return redis.error_reply("Invalid user.")	
end

-- check if user owns any groups
if redis.call("HGET", userHashKey, "adminGroupCount") ~= "0" then
	return redis.error_reply("Admin group count not 0.")	
end

-- get users email from user hash
local userEmail = redis.call("HGET", userHashKey, "email")

-- key deletes
redis.call("DEL", userHashKey)
redis.call("DEL", userGroupsKey)
redis.call("DEL", userGroupInvitesKey)
redis.call("HDEL", userIDKey, userEmail)

-- get groupID's from user group list
local userGroupList = redis.call("HKEYS", userGroupsKey)

for i = 1, #userGroupList do
	-- delete user group messages
	redis.call("DEL", userGrpMsgKeyPartial .. userID .. ":" .. userGroupList[i])
	-- delete user from group
	redis.call("HDEL", groupMemKeyPartial .. userGroupList[i], userID)
	-- delete user from all group locations
	redis.call("HDEL", groupLocationsKeyPartial .. userGroupList[i], userID)
end
--

return "Success"
