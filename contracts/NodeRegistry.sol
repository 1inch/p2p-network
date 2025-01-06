// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract NodeRegistry {
    struct Resolver {
        string ip;
    }

    string private relayerIP;

    mapping(bytes => Resolver) private resolvers;

    bytes[] private resolverKeys;

    /// @notice Register the relayer node with its IP
    /// @param ip The IP address of the relayer node
    function registerRelayer(string calldata ip) external {
        require(bytes(ip).length > 0, "Relayer IP cannot be empty");
        relayerIP = ip;
    }

    /// @notice Register a resolver node with its IP and public key
    /// @param ip The IP address of the resolver node
    /// @param publicKey The public key of the resolver node as bytes
    function registerResolver(string calldata ip, bytes calldata publicKey) external {
        require(bytes(ip).length > 0, "Resolver IP cannot be empty");
        require(publicKey.length > 0, "Public key cannot be empty");
        require(bytes(resolvers[publicKey].ip).length == 0, "Resolver already registered");

        resolvers[publicKey] = Resolver({
            ip: ip
        });

        resolverKeys.push(publicKey);
    }

    /// @notice Get the IP address of the relayer node and all resolver public keys
    /// @return ip The IP address of the relayer node
    /// @return publicKeys An array of all resolver public keys
    function getRelayer() external view returns (string memory ip, bytes[] memory publicKeys) {
        require(bytes(relayerIP).length > 0, "No relayer registered");
        ip = relayerIP;
        publicKeys = resolverKeys;
    }

    /// @notice Get the IP address of a resolver node by its public key
    /// @param publicKey The public key of the resolver node as bytes
    /// @return ip The IP address of the resolver node
    function getResolver(bytes calldata publicKey) external view returns (string memory ip) {
        string memory resolverIP = resolvers[publicKey].ip;
        require(bytes(resolverIP).length > 0, "Resolver not found");
        return resolverIP;
    }
}
