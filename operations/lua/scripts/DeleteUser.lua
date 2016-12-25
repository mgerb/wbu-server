local userID = ARGV[1]
local userGrpMsgKey = ARGV[2]
local groupMemKey = ARGV[3]
local userHashKey = KEYS[1]
local userIDKey = KEYS[2]
local userGroupsKey = KEYS[3]
local userGroupInvitesKey = KEYS[4]

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
	redis.call("DEL", userGrpMsgKey .. userID .. userGroupList[i])
	-- delete user from group
	redis.call("HDEL", groupMemKey .. userGroupList[i], userID)
end
--

return "Success"
