local groupHashKey = KEYS[1]
local groupMembersKey = KEYS[2]
local groupLocationsKey = KEYS[3]
local userGrpMsgKey = KEYS[4]
local userGroupsKey = KEYS[5]

local userID = ARGV[1]
local ownerID = ARGV[2]
local groupID = ARGV[3]

-- check if owner has permissions
if redis.call("HGET", groupHashKey, "owner") ~= ownerID then
    return redis.error_reply("Invalid permissions.")
end

-- check if user exists in group
if redis.call("HEXISTS", groupMemKey, userID) == 0 then
	return redis.error_reply("User doesn't exist in group.")	
end

-- perform deletions
redis.call("HDEL", groupMembersKey, userID)
redis.call("HDEL", groupLocationsKey, userID)
redis.call("HDEL", userGroupsKey, groupID)
redis.call("DEL", userGrpMsgKey)

return "Success"
