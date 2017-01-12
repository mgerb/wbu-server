-- gets users in a group along with the owner
local groupHashKey = KEYS[1]
local groupMembersKey = KEYS[2]

local userID = ARGV[1]

-- check if user exists in group
if redis.call("HEXISTS", groupMembersKey, userID) == 0 then
    return redis.error_reply("Invalid permissions")
end

-- get group owner
local groupOwner = redis.call("HGET", groupHashKey, "owner")

-- get all group members
local groupMembers = redis.call("HGETALL", groupMembersKey)

local resultSet = {}

for i = 1, #groupMembers, 2 do

    resultSet[#resultSet+1] = { 
        userID = groupMembers[i],
        fullName = groupMembers[i+1],
        owner = groupOwner == groupMembers[i]
    }

end

return cjson.encode(resultSet)
