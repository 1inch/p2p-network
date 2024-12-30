pragma solidity ^0.8.0;

contract NodeRegistry {
    struct Resolver {
        string ip;
        bytes publicKey;
        bool exists;
    }

    string private relayerIP;
    bool private relayerExists;

    mapping(bytes => Resolver) private resolvers;

    event RelayerRegistered(string ip);
    event ResolverRegistered(string ip, bytes publicKey);

    /// @notice Register the relayer node with its IP
    /// @param ip The IP address of the relayer node
    function registerRelayer(string calldata ip) external {
        relayerIP = ip;
        relayerExists = true;

        emit RelayerRegistered(ip);
    }

    /// @notice Register a resolver node with its IP and public key
    /// @param ip The IP address of the resolver node
    /// @param publicKey The public key of the resolver node as bytes
    function registerResolver(string calldata ip, bytes calldata publicKey) external {
        require(!resolvers[publicKey].exists, "Resolver already registered");

        resolvers[publicKey] = Resolver({
            ip: ip,
            publicKey: publicKey,
            exists: true
        });

        emit ResolverRegistered(ip, publicKey);
    }

    /// @notice Get details of the current relayer node
    /// @return ip The IP address of the relayer node
    function getRelayer() external view returns (string memory ip) {
        require(relayerExists, "No relayer registered");
        return relayerIP;
    }

    /// @notice Get details of a resolver node by its public key
    /// @param publicKey The public key of the resolver node as bytes
    /// @return ip The IP address of the resolver node
    function getResolver(bytes calldata publicKey) external view returns (string memory ip) {
        require(resolvers[publicKey].exists, "Resolver not found");
        return resolvers[publicKey].ip;
    }
}
