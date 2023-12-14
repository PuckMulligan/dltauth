# dltauth

Concept Overview
Dual Hashing Mechanism:

When a user inputs a password, it is processed through two separate hashing algorithms.
First Hash (Local Authentication): One hash is used for the traditional authentication process on the local system.
Second Hash (Blockchain Evidence): The other hash, generated using a different algorithm, is recorded on the blockchain.
Purpose:

The first hash serves the usual purpose of password verification.
The second hash acts as evidence of the actual password entry, aiming to differentiate between a genuine user input and a stolen hash used in a PtH attack. This is a demo and as such there are concerns that need to be addressed

- Blockchain introduces significant overhead. Hyperledger Fabric was chosen for its ease of use and readiness for integration for the purposes of a demo. However, in a real world scenario, a distributed ledger tailored to this purpose would be constructed. This is currently beyond the scope of this project.