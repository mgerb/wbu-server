local groupIDKey = KEYS[1]
local userIDKey = KEYS[2]

local userGrpMsgKeyPartial = ARGV[1]
local userID = ARGV[2]
local groupID = ARGV[3]
local timeStamp = ARGV[4]
local message = ARGV[5]

local oneMonth = 2592000


local fullName = redis.call("HGET", userIDKey, "fullName")

-- get user full name
if not fullName then
	return redis.error_reply("User does not exist")	
end

--check if user exists in group
if redis.call("HEXISTS", groupIDKey, userID) == 0 then
	return redis.error_reply("Invalid permissions.")
end

local fullMessage = userID .. "/" .. fullName .. "/" .. timeStamp .. "/" .. message

-- get all the members in the group
local members = redis.call("HKEYS", groupIDKey)

-- cycle through each key in the group member hash
for i = 1, #members do
	redis.call("SADD", userGrpMsgKeyPartial .. members[i] .. ":" .. groupID, fullMessage)
	-- reset the expire time for each message
	redis.call("EXPIRE", userGrpMsgKeyPartial .. members[i] .. ":" .. groupID, oneMonth)
end

return "Success"