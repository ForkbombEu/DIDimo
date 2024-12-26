import PocketBase from 'pocketbase';
import { zencode_exec } from "zenroom";
const KEYPAIROOM_CLIENT = `
Scenario 'ecdh': Create the key
Scenario 'ethereum': Create key
Scenario 'reflow': Create the key
Scenario 'eddsa': Create the key


# Loading the user name from data
Given my name is in a 'string' named 'username'

# Loading the answers from 3 secret questions. The user will have to pick the 3 challenges from a list
# and have to remember the questions - the order is not important cause Zenroom will sort alphabetically
# the data in input
#
# NOTE: the challenges will never be communicated to the server or to anybody else!
Given I have a 'string dictionary' named 'userChallenges'

# Loading the individual challenges, in order to have them hashed
# and the hashes OPTIONALLY stored by the server, to improve regeneration of the keypair
Given I have a 'string' named 'whereParentsMet' in 'userChallenges'
Given I have a 'string' named 'nameFirstPet' in 'userChallenges'
Given I have a 'string' named 'whereHomeTown' in 'userChallenges'
Given I have a 'string' named 'nameFirstTeacher' in 'userChallenges'
Given I have a 'string' named 'nameMotherMaid' in 'userChallenges'

# Loading the pbkdf received from the server, containing a signed hash of known data
Given that I have a 'base64' named 'seedServerSideShard.HMAC'

# Save the backup for mnemonic dump, before factoring with the salt
# it is shortened to 16 bytes by hashing sha512 the KDF and taking the first 16 bytes
When I create the key derivation of 'userChallenges'
and I create the hash of 'key derivation' using 'sha512'
and I split the leftmost '16' bytes of 'hash'
and I delete the 'key derivation'
and I delete the 'hash'
and I rename the 'leftmost' to 'seed'

# Hash again the user's challenges with salt for the seed root
When I rename 'seedServerSideShard.HMAC' to 'salt'
and I create the key derivation of 'seed' with password 'salt'
and I rename the 'key derivation' to 'seed.root'

# In the following flow the order should NOT be changed

When I create the hash of 'seed.root'
When I rename the 'hash' to 'seed.ecdh'

When I create the hash of 'seed.ecdh'
When I rename the 'hash' to 'seed.eddsa'

When I create the hash of 'seed.eddsa'
When I rename the 'hash' to 'seed.ethereum'

When I create the hash of 'seed.ethereum'
When I rename the 'hash' to 'seed.reflow'

When I create the hash of 'seed.reflow'
When I rename the 'hash' to 'seed.bitcoin'

# end of the sorted creation flow

When I create the ecdh key with secret key 'seed.ecdh'
When I create the eddsa key with secret key 'seed.eddsa'
When I create the ethereum key with secret key 'seed.ethereum'
When I create the reflow key with secret key 'seed.reflow'
When I create the bitcoin key with secret key 'seed.bitcoin'

When I create the ecdh public key
When I create the eddsa public key
When I create the ethereum address
When I create the reflow public key
When I create the bitcoin public key

# Creating the hashes of the single challenges, to OPTIONALLY help
# regeneration of the keypair

# When I create the 'base64 dictionary'
# and I rename the 'base64 dictionary' to 'hashedAnswers'

# When I create the key derivation of 'whereParentsMet'
# and I rename the 'key derivation' to 'whereParentsMet.kdf'
# When I move 'whereParentsMet.kdf' in 'hashedAnswers'

# When I create the key derivation of 'nameFirstPet'
# and I rename the 'key derivation' to 'nameFirstPet.kdf'
# When I move 'nameFirstPet.kdf' in 'hashedAnswers'

# When I create the key derivation of 'whereHomeTown'
# and I rename the 'key derivation' to 'whereHomeTown.kdf'
# When I move 'whereHomeTown.kdf' in 'hashedAnswers'

# When I create the key derivation of 'nameFirstTeacher'
# and I rename the 'key derivation' to 'nameFirstTeacher.kdf'
# When I move 'nameFirstTeacher.kdf' in 'hashedAnswers'

# When I create the key derivation of 'nameMotherMaid'
# and I rename the 'key derivation' to 'nameMotherMaid.kdf'
# When I move 'nameMotherMaid.kdf' in 'hashedAnswers'

# this prints the hashes of the challenges
# Then print the 'hashedAnswers'

Then print the 'keyring'

# This prints the seed for the private keys as mnemonic
Then print the 'seed' as 'mnemonic'

Then print the 'ecdh public key' as 'base58'
Then print the 'eddsa public key' as 'base58'
Then print the 'ethereum address'
Then print the 'reflow public key' as 'base58'
Then print the 'bitcoin public key' as 'base58'
`
const pb = new PocketBase('https://test.signroom.io');
const email = "mimmo@sito.io"

const hmac = await pb.send("/api/keypairoom-server", {
    body: JSON.stringify({ "userData": {"email": email}}), method: "POST", headers: {
      "Content-Type": "application/json",}});
console.log(hmac)
const zenData = {
    userChallenges: {
        "whereParentsMet": "mimmo",
        "nameFirstPet": "mimmo",
        "nameFirstTeacher": "mimmo",
        "whereHomeTown": "mimmo",
        "nameMotherMaid": "mimmo",
    },
    username: email,
    "seedServerSideShard.HMAC": hmac.hmac
};

const { result } = await zencode_exec(KEYPAIROOM_CLIENT, { data: JSON.stringify(zenData) });
console.log(result)

