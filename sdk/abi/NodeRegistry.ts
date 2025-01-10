export const registryAbi = [
    {
        inputs: [],
        name: "getRelayer",
        outputs: [
            {
                internalType: "string",
                name: "ip",
                type: "string"
            },
            {
                internalType: "bytes[]",
                name: "publicKeys",
                type: "bytes[]"
            }
        ],
        stateMutability: "view",
        type: "function"
    },
    {
        inputs: [
            {
                internalType: "bytes",
                name: "publicKey",
                type: "bytes"
            }
        ],
        name: "getResolver",
        outputs: [
            {
                internalType: "string",
                name: "ip",
                type: "string"
            }
        ],
        stateMutability: "view",
        type: "function"
    },
    {
        inputs: [
            {
                internalType: "string",
                name: "ip",
                yype: "string"
            }
        ],
        name: "registerRelayer",
        outputs: [],
        stateMutability: "nonpayable",
        type: "function"
    },
    {
        inputs: [
            {
                internalType: "string",
                name: "ip",
                type: "string"
            },
            {
                internalType: "bytes",
                name: "publicKey",
                type: "bytes"
            }
        ],
        name: "registerResolver",
        outputs: [],
        stateMutability: "nonpayable",
        type: "function"
    }
];
