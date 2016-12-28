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

local groupMembers = redis.call("HGETALL", groupMembersKey)

local t = {}
local userIsOwner = false

for i = 2, #groupMembers, 2 do

    if groupOwner == groupMembers[i] then
        userIsOwner = true
    end

    t[#t+1] = { 
        userID = groupMembers[i],
        fullName = groupMembers[i+1],
        owner = userIsOwner
    }

    userIsOwner = false

end

return cjson.encode({ owner = groupOwner, members = t })
