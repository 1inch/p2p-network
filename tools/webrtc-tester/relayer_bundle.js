/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
(function($protobuf) {
    "use strict";

    // Common aliases
    var $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;
    
    // Exported root namespace
    var $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});
    
    $root.proto = (function() {
    
        /**
         * Namespace proto.
         * @exports proto
         * @namespace
         */
        var proto = {};
    
        /**
         * ErrorCode enum.
         * @name proto.ErrorCode
         * @enum {number}
         * @property {number} ERR_INVALID_MESSAGE_FORMAT=0 ERR_INVALID_MESSAGE_FORMAT value
         * @property {number} ERR_RESOLVER_LOOKUP_FAILED=1 ERR_RESOLVER_LOOKUP_FAILED value
         * @property {number} ERR_GRPC_EXECUTION_FAILED=2 ERR_GRPC_EXECUTION_FAILED value
         * @property {number} ERR_RESPONSE_SERIALIZATION_FAILED=3 ERR_RESPONSE_SERIALIZATION_FAILED value
         * @property {number} ERR_DATA_CHANNEL_SEND_FAILED=4 ERR_DATA_CHANNEL_SEND_FAILED value
         */
        proto.ErrorCode = (function() {
            var valuesById = {}, values = Object.create(valuesById);
            values[valuesById[0] = "ERR_INVALID_MESSAGE_FORMAT"] = 0;
            values[valuesById[1] = "ERR_RESOLVER_LOOKUP_FAILED"] = 1;
            values[valuesById[2] = "ERR_GRPC_EXECUTION_FAILED"] = 2;
            values[valuesById[3] = "ERR_RESPONSE_SERIALIZATION_FAILED"] = 3;
            values[valuesById[4] = "ERR_DATA_CHANNEL_SEND_FAILED"] = 4;
            return values;
        })();
    
        proto.Error = (function() {
    
            /**
             * Properties of an Error.
             * @memberof proto
             * @interface IError
             * @property {proto.ErrorCode|null} [code] Error code
             * @property {string|null} [message] Error message
             */
    
            /**
             * Constructs a new Error.
             * @memberof proto
             * @classdesc Represents an Error.
             * @implements IError
             * @constructor
             * @param {proto.IError=} [properties] Properties to set
             */
            function Error(properties) {
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }
    
            /**
             * Error code.
             * @member {proto.ErrorCode} code
             * @memberof proto.Error
             * @instance
             */
            Error.prototype.code = 0;
    
            /**
             * Error message.
             * @member {string} message
             * @memberof proto.Error
             * @instance
             */
            Error.prototype.message = "";
    
            /**
             * Creates a new Error instance using the specified properties.
             * @function create
             * @memberof proto.Error
             * @static
             * @param {proto.IError=} [properties] Properties to set
             * @returns {proto.Error} Error instance
             */
            Error.create = function create(properties) {
                return new Error(properties);
            };
    
            /**
             * Encodes the specified Error message. Does not implicitly {@link proto.Error.verify|verify} messages.
             * @function encode
             * @memberof proto.Error
             * @static
             * @param {proto.IError} message Error message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Error.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.code != null && Object.hasOwnProperty.call(message, "code"))
                    writer.uint32(/* id 1, wireType 0 =*/8).int32(message.code);
                if (message.message != null && Object.hasOwnProperty.call(message, "message"))
                    writer.uint32(/* id 2, wireType 2 =*/18).string(message.message);
                return writer;
            };
    
            /**
             * Encodes the specified Error message, length delimited. Does not implicitly {@link proto.Error.verify|verify} messages.
             * @function encodeDelimited
             * @memberof proto.Error
             * @static
             * @param {proto.IError} message Error message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Error.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };
    
            /**
             * Decodes an Error message from the specified reader or buffer.
             * @function decode
             * @memberof proto.Error
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {proto.Error} Error
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Error.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.proto.Error();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1: {
                            message.code = reader.int32();
                            break;
                        }
                    case 2: {
                            message.message = reader.string();
                            break;
                        }
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };
    
            /**
             * Decodes an Error message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof proto.Error
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {proto.Error} Error
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Error.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };
    
            /**
             * Verifies an Error message.
             * @function verify
             * @memberof proto.Error
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            Error.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.code != null && message.hasOwnProperty("code"))
                    switch (message.code) {
                    default:
                        return "code: enum value expected";
                    case 0:
                    case 1:
                    case 2:
                    case 3:
                    case 4:
                        break;
                    }
                if (message.message != null && message.hasOwnProperty("message"))
                    if (!$util.isString(message.message))
                        return "message: string expected";
                return null;
            };
    
            /**
             * Creates an Error message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof proto.Error
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {proto.Error} Error
             */
            Error.fromObject = function fromObject(object) {
                if (object instanceof $root.proto.Error)
                    return object;
                var message = new $root.proto.Error();
                switch (object.code) {
                default:
                    if (typeof object.code === "number") {
                        message.code = object.code;
                        break;
                    }
                    break;
                case "ERR_INVALID_MESSAGE_FORMAT":
                case 0:
                    message.code = 0;
                    break;
                case "ERR_RESOLVER_LOOKUP_FAILED":
                case 1:
                    message.code = 1;
                    break;
                case "ERR_GRPC_EXECUTION_FAILED":
                case 2:
                    message.code = 2;
                    break;
                case "ERR_RESPONSE_SERIALIZATION_FAILED":
                case 3:
                    message.code = 3;
                    break;
                case "ERR_DATA_CHANNEL_SEND_FAILED":
                case 4:
                    message.code = 4;
                    break;
                }
                if (object.message != null)
                    message.message = String(object.message);
                return message;
            };
    
            /**
             * Creates a plain object from an Error message. Also converts values to other types if specified.
             * @function toObject
             * @memberof proto.Error
             * @static
             * @param {proto.Error} message Error
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            Error.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.defaults) {
                    object.code = options.enums === String ? "ERR_INVALID_MESSAGE_FORMAT" : 0;
                    object.message = "";
                }
                if (message.code != null && message.hasOwnProperty("code"))
                    object.code = options.enums === String ? $root.proto.ErrorCode[message.code] === undefined ? message.code : $root.proto.ErrorCode[message.code] : message.code;
                if (message.message != null && message.hasOwnProperty("message"))
                    object.message = message.message;
                return object;
            };
    
            /**
             * Converts this Error to JSON.
             * @function toJSON
             * @memberof proto.Error
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            Error.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };
    
            /**
             * Gets the default type url for Error
             * @function getTypeUrl
             * @memberof proto.Error
             * @static
             * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns {string} The default type url
             */
            Error.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
                if (typeUrlPrefix === undefined) {
                    typeUrlPrefix = "type.googleapis.com";
                }
                return typeUrlPrefix + "/proto.Error";
            };
    
            return Error;
        })();
    
        proto.IncomingMessage = (function() {
    
            /**
             * Properties of an IncomingMessage.
             * @memberof proto
             * @interface IIncomingMessage
             * @property {Array.<Uint8Array>|null} [publicKeys] IncomingMessage publicKeys
             * @property {proto.IResolverRequest|null} [request] IncomingMessage request
             */
    
            /**
             * Constructs a new IncomingMessage.
             * @memberof proto
             * @classdesc Represents an IncomingMessage.
             * @implements IIncomingMessage
             * @constructor
             * @param {proto.IIncomingMessage=} [properties] Properties to set
             */
            function IncomingMessage(properties) {
                this.publicKeys = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }
    
            /**
             * IncomingMessage publicKeys.
             * @member {Array.<Uint8Array>} publicKeys
             * @memberof proto.IncomingMessage
             * @instance
             */
            IncomingMessage.prototype.publicKeys = $util.emptyArray;
    
            /**
             * IncomingMessage request.
             * @member {proto.IResolverRequest|null|undefined} request
             * @memberof proto.IncomingMessage
             * @instance
             */
            IncomingMessage.prototype.request = null;
    
            /**
             * Creates a new IncomingMessage instance using the specified properties.
             * @function create
             * @memberof proto.IncomingMessage
             * @static
             * @param {proto.IIncomingMessage=} [properties] Properties to set
             * @returns {proto.IncomingMessage} IncomingMessage instance
             */
            IncomingMessage.create = function create(properties) {
                return new IncomingMessage(properties);
            };
    
            /**
             * Encodes the specified IncomingMessage message. Does not implicitly {@link proto.IncomingMessage.verify|verify} messages.
             * @function encode
             * @memberof proto.IncomingMessage
             * @static
             * @param {proto.IIncomingMessage} message IncomingMessage message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            IncomingMessage.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.publicKeys != null && message.publicKeys.length)
                    for (var i = 0; i < message.publicKeys.length; ++i)
                        writer.uint32(/* id 1, wireType 2 =*/10).bytes(message.publicKeys[i]);
                if (message.request != null && Object.hasOwnProperty.call(message, "request"))
                    $root.proto.ResolverRequest.encode(message.request, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
                return writer;
            };
    
            /**
             * Encodes the specified IncomingMessage message, length delimited. Does not implicitly {@link proto.IncomingMessage.verify|verify} messages.
             * @function encodeDelimited
             * @memberof proto.IncomingMessage
             * @static
             * @param {proto.IIncomingMessage} message IncomingMessage message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            IncomingMessage.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };
    
            /**
             * Decodes an IncomingMessage message from the specified reader or buffer.
             * @function decode
             * @memberof proto.IncomingMessage
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {proto.IncomingMessage} IncomingMessage
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            IncomingMessage.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.proto.IncomingMessage();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1: {
                            if (!(message.publicKeys && message.publicKeys.length))
                                message.publicKeys = [];
                            message.publicKeys.push(reader.bytes());
                            break;
                        }
                    case 2: {
                            message.request = $root.proto.ResolverRequest.decode(reader, reader.uint32());
                            break;
                        }
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };
    
            /**
             * Decodes an IncomingMessage message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof proto.IncomingMessage
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {proto.IncomingMessage} IncomingMessage
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            IncomingMessage.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };
    
            /**
             * Verifies an IncomingMessage message.
             * @function verify
             * @memberof proto.IncomingMessage
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            IncomingMessage.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.publicKeys != null && message.hasOwnProperty("publicKeys")) {
                    if (!Array.isArray(message.publicKeys))
                        return "publicKeys: array expected";
                    for (var i = 0; i < message.publicKeys.length; ++i)
                        if (!(message.publicKeys[i] && typeof message.publicKeys[i].length === "number" || $util.isString(message.publicKeys[i])))
                            return "publicKeys: buffer[] expected";
                }
                if (message.request != null && message.hasOwnProperty("request")) {
                    var error = $root.proto.ResolverRequest.verify(message.request);
                    if (error)
                        return "request." + error;
                }
                return null;
            };
    
            /**
             * Creates an IncomingMessage message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof proto.IncomingMessage
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {proto.IncomingMessage} IncomingMessage
             */
            IncomingMessage.fromObject = function fromObject(object) {
                if (object instanceof $root.proto.IncomingMessage)
                    return object;
                var message = new $root.proto.IncomingMessage();
                if (object.publicKeys) {
                    if (!Array.isArray(object.publicKeys))
                        throw TypeError(".proto.IncomingMessage.publicKeys: array expected");
                    message.publicKeys = [];
                    for (var i = 0; i < object.publicKeys.length; ++i)
                        if (typeof object.publicKeys[i] === "string")
                            $util.base64.decode(object.publicKeys[i], message.publicKeys[i] = $util.newBuffer($util.base64.length(object.publicKeys[i])), 0);
                        else if (object.publicKeys[i].length >= 0)
                            message.publicKeys[i] = object.publicKeys[i];
                }
                if (object.request != null) {
                    if (typeof object.request !== "object")
                        throw TypeError(".proto.IncomingMessage.request: object expected");
                    message.request = $root.proto.ResolverRequest.fromObject(object.request);
                }
                return message;
            };
    
            /**
             * Creates a plain object from an IncomingMessage message. Also converts values to other types if specified.
             * @function toObject
             * @memberof proto.IncomingMessage
             * @static
             * @param {proto.IncomingMessage} message IncomingMessage
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            IncomingMessage.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults)
                    object.publicKeys = [];
                if (options.defaults)
                    object.request = null;
                if (message.publicKeys && message.publicKeys.length) {
                    object.publicKeys = [];
                    for (var j = 0; j < message.publicKeys.length; ++j)
                        object.publicKeys[j] = options.bytes === String ? $util.base64.encode(message.publicKeys[j], 0, message.publicKeys[j].length) : options.bytes === Array ? Array.prototype.slice.call(message.publicKeys[j]) : message.publicKeys[j];
                }
                if (message.request != null && message.hasOwnProperty("request"))
                    object.request = $root.proto.ResolverRequest.toObject(message.request, options);
                return object;
            };
    
            /**
             * Converts this IncomingMessage to JSON.
             * @function toJSON
             * @memberof proto.IncomingMessage
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            IncomingMessage.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };
    
            /**
             * Gets the default type url for IncomingMessage
             * @function getTypeUrl
             * @memberof proto.IncomingMessage
             * @static
             * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns {string} The default type url
             */
            IncomingMessage.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
                if (typeUrlPrefix === undefined) {
                    typeUrlPrefix = "type.googleapis.com";
                }
                return typeUrlPrefix + "/proto.IncomingMessage";
            };
    
            return IncomingMessage;
        })();
    
        proto.OutgoingMessage = (function() {
    
            /**
             * Properties of an OutgoingMessage.
             * @memberof proto
             * @interface IOutgoingMessage
             * @property {proto.IResolverResponse|null} [response] OutgoingMessage response
             * @property {proto.IError|null} [error] OutgoingMessage error
             * @property {Uint8Array|null} [publicKey] OutgoingMessage publicKey
             */
    
            /**
             * Constructs a new OutgoingMessage.
             * @memberof proto
             * @classdesc Represents an OutgoingMessage.
             * @implements IOutgoingMessage
             * @constructor
             * @param {proto.IOutgoingMessage=} [properties] Properties to set
             */
            function OutgoingMessage(properties) {
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }
    
            /**
             * OutgoingMessage response.
             * @member {proto.IResolverResponse|null|undefined} response
             * @memberof proto.OutgoingMessage
             * @instance
             */
            OutgoingMessage.prototype.response = null;
    
            /**
             * OutgoingMessage error.
             * @member {proto.IError|null|undefined} error
             * @memberof proto.OutgoingMessage
             * @instance
             */
            OutgoingMessage.prototype.error = null;
    
            /**
             * OutgoingMessage publicKey.
             * @member {Uint8Array} publicKey
             * @memberof proto.OutgoingMessage
             * @instance
             */
            OutgoingMessage.prototype.publicKey = $util.newBuffer([]);
    
            // OneOf field names bound to virtual getters and setters
            var $oneOfFields;
    
            /**
             * OutgoingMessage result.
             * @member {"response"|"error"|undefined} result
             * @memberof proto.OutgoingMessage
             * @instance
             */
            Object.defineProperty(OutgoingMessage.prototype, "result", {
                get: $util.oneOfGetter($oneOfFields = ["response", "error"]),
                set: $util.oneOfSetter($oneOfFields)
            });
    
            /**
             * Creates a new OutgoingMessage instance using the specified properties.
             * @function create
             * @memberof proto.OutgoingMessage
             * @static
             * @param {proto.IOutgoingMessage=} [properties] Properties to set
             * @returns {proto.OutgoingMessage} OutgoingMessage instance
             */
            OutgoingMessage.create = function create(properties) {
                return new OutgoingMessage(properties);
            };
    
            /**
             * Encodes the specified OutgoingMessage message. Does not implicitly {@link proto.OutgoingMessage.verify|verify} messages.
             * @function encode
             * @memberof proto.OutgoingMessage
             * @static
             * @param {proto.IOutgoingMessage} message OutgoingMessage message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            OutgoingMessage.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.response != null && Object.hasOwnProperty.call(message, "response"))
                    $root.proto.ResolverResponse.encode(message.response, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
                if (message.error != null && Object.hasOwnProperty.call(message, "error"))
                    $root.proto.Error.encode(message.error, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
                if (message.publicKey != null && Object.hasOwnProperty.call(message, "publicKey"))
                    writer.uint32(/* id 3, wireType 2 =*/26).bytes(message.publicKey);
                return writer;
            };
    
            /**
             * Encodes the specified OutgoingMessage message, length delimited. Does not implicitly {@link proto.OutgoingMessage.verify|verify} messages.
             * @function encodeDelimited
             * @memberof proto.OutgoingMessage
             * @static
             * @param {proto.IOutgoingMessage} message OutgoingMessage message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            OutgoingMessage.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };
    
            /**
             * Decodes an OutgoingMessage message from the specified reader or buffer.
             * @function decode
             * @memberof proto.OutgoingMessage
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {proto.OutgoingMessage} OutgoingMessage
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            OutgoingMessage.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.proto.OutgoingMessage();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1: {
                            message.response = $root.proto.ResolverResponse.decode(reader, reader.uint32());
                            break;
                        }
                    case 2: {
                            message.error = $root.proto.Error.decode(reader, reader.uint32());
                            break;
                        }
                    case 3: {
                            message.publicKey = reader.bytes();
                            break;
                        }
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };
    
            /**
             * Decodes an OutgoingMessage message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof proto.OutgoingMessage
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {proto.OutgoingMessage} OutgoingMessage
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            OutgoingMessage.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };
    
            /**
             * Verifies an OutgoingMessage message.
             * @function verify
             * @memberof proto.OutgoingMessage
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            OutgoingMessage.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                var properties = {};
                if (message.response != null && message.hasOwnProperty("response")) {
                    properties.result = 1;
                    {
                        var error = $root.proto.ResolverResponse.verify(message.response);
                        if (error)
                            return "response." + error;
                    }
                }
                if (message.error != null && message.hasOwnProperty("error")) {
                    if (properties.result === 1)
                        return "result: multiple values";
                    properties.result = 1;
                    {
                        var error = $root.proto.Error.verify(message.error);
                        if (error)
                            return "error." + error;
                    }
                }
                if (message.publicKey != null && message.hasOwnProperty("publicKey"))
                    if (!(message.publicKey && typeof message.publicKey.length === "number" || $util.isString(message.publicKey)))
                        return "publicKey: buffer expected";
                return null;
            };
    
            /**
             * Creates an OutgoingMessage message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof proto.OutgoingMessage
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {proto.OutgoingMessage} OutgoingMessage
             */
            OutgoingMessage.fromObject = function fromObject(object) {
                if (object instanceof $root.proto.OutgoingMessage)
                    return object;
                var message = new $root.proto.OutgoingMessage();
                if (object.response != null) {
                    if (typeof object.response !== "object")
                        throw TypeError(".proto.OutgoingMessage.response: object expected");
                    message.response = $root.proto.ResolverResponse.fromObject(object.response);
                }
                if (object.error != null) {
                    if (typeof object.error !== "object")
                        throw TypeError(".proto.OutgoingMessage.error: object expected");
                    message.error = $root.proto.Error.fromObject(object.error);
                }
                if (object.publicKey != null)
                    if (typeof object.publicKey === "string")
                        $util.base64.decode(object.publicKey, message.publicKey = $util.newBuffer($util.base64.length(object.publicKey)), 0);
                    else if (object.publicKey.length >= 0)
                        message.publicKey = object.publicKey;
                return message;
            };
    
            /**
             * Creates a plain object from an OutgoingMessage message. Also converts values to other types if specified.
             * @function toObject
             * @memberof proto.OutgoingMessage
             * @static
             * @param {proto.OutgoingMessage} message OutgoingMessage
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            OutgoingMessage.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.defaults)
                    if (options.bytes === String)
                        object.publicKey = "";
                    else {
                        object.publicKey = [];
                        if (options.bytes !== Array)
                            object.publicKey = $util.newBuffer(object.publicKey);
                    }
                if (message.response != null && message.hasOwnProperty("response")) {
                    object.response = $root.proto.ResolverResponse.toObject(message.response, options);
                    if (options.oneofs)
                        object.result = "response";
                }
                if (message.error != null && message.hasOwnProperty("error")) {
                    object.error = $root.proto.Error.toObject(message.error, options);
                    if (options.oneofs)
                        object.result = "error";
                }
                if (message.publicKey != null && message.hasOwnProperty("publicKey"))
                    object.publicKey = options.bytes === String ? $util.base64.encode(message.publicKey, 0, message.publicKey.length) : options.bytes === Array ? Array.prototype.slice.call(message.publicKey) : message.publicKey;
                return object;
            };
    
            /**
             * Converts this OutgoingMessage to JSON.
             * @function toJSON
             * @memberof proto.OutgoingMessage
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            OutgoingMessage.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };
    
            /**
             * Gets the default type url for OutgoingMessage
             * @function getTypeUrl
             * @memberof proto.OutgoingMessage
             * @static
             * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns {string} The default type url
             */
            OutgoingMessage.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
                if (typeUrlPrefix === undefined) {
                    typeUrlPrefix = "type.googleapis.com";
                }
                return typeUrlPrefix + "/proto.OutgoingMessage";
            };
    
            return OutgoingMessage;
        })();
    
        proto.ResolverRequest = (function() {
    
            /**
             * Properties of a ResolverRequest.
             * @memberof proto
             * @interface IResolverRequest
             * @property {string|null} [id] ResolverRequest id
             * @property {Uint8Array|null} [payload] ResolverRequest payload
             * @property {boolean|null} [encrypted] ResolverRequest encrypted
             * @property {Uint8Array|null} [publicKey] ResolverRequest publicKey
             */
    
            /**
             * Constructs a new ResolverRequest.
             * @memberof proto
             * @classdesc Represents a ResolverRequest.
             * @implements IResolverRequest
             * @constructor
             * @param {proto.IResolverRequest=} [properties] Properties to set
             */
            function ResolverRequest(properties) {
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }
    
            /**
             * ResolverRequest id.
             * @member {string} id
             * @memberof proto.ResolverRequest
             * @instance
             */
            ResolverRequest.prototype.id = "";
    
            /**
             * ResolverRequest payload.
             * @member {Uint8Array} payload
             * @memberof proto.ResolverRequest
             * @instance
             */
            ResolverRequest.prototype.payload = $util.newBuffer([]);
    
            /**
             * ResolverRequest encrypted.
             * @member {boolean} encrypted
             * @memberof proto.ResolverRequest
             * @instance
             */
            ResolverRequest.prototype.encrypted = false;
    
            /**
             * ResolverRequest publicKey.
             * @member {Uint8Array} publicKey
             * @memberof proto.ResolverRequest
             * @instance
             */
            ResolverRequest.prototype.publicKey = $util.newBuffer([]);
    
            /**
             * Creates a new ResolverRequest instance using the specified properties.
             * @function create
             * @memberof proto.ResolverRequest
             * @static
             * @param {proto.IResolverRequest=} [properties] Properties to set
             * @returns {proto.ResolverRequest} ResolverRequest instance
             */
            ResolverRequest.create = function create(properties) {
                return new ResolverRequest(properties);
            };
    
            /**
             * Encodes the specified ResolverRequest message. Does not implicitly {@link proto.ResolverRequest.verify|verify} messages.
             * @function encode
             * @memberof proto.ResolverRequest
             * @static
             * @param {proto.IResolverRequest} message ResolverRequest message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            ResolverRequest.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
                if (message.payload != null && Object.hasOwnProperty.call(message, "payload"))
                    writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.payload);
                if (message.encrypted != null && Object.hasOwnProperty.call(message, "encrypted"))
                    writer.uint32(/* id 3, wireType 0 =*/24).bool(message.encrypted);
                if (message.publicKey != null && Object.hasOwnProperty.call(message, "publicKey"))
                    writer.uint32(/* id 4, wireType 2 =*/34).bytes(message.publicKey);
                return writer;
            };
    
            /**
             * Encodes the specified ResolverRequest message, length delimited. Does not implicitly {@link proto.ResolverRequest.verify|verify} messages.
             * @function encodeDelimited
             * @memberof proto.ResolverRequest
             * @static
             * @param {proto.IResolverRequest} message ResolverRequest message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            ResolverRequest.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };
    
            /**
             * Decodes a ResolverRequest message from the specified reader or buffer.
             * @function decode
             * @memberof proto.ResolverRequest
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {proto.ResolverRequest} ResolverRequest
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            ResolverRequest.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.proto.ResolverRequest();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1: {
                            message.id = reader.string();
                            break;
                        }
                    case 2: {
                            message.payload = reader.bytes();
                            break;
                        }
                    case 3: {
                            message.encrypted = reader.bool();
                            break;
                        }
                    case 4: {
                            message.publicKey = reader.bytes();
                            break;
                        }
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };
    
            /**
             * Decodes a ResolverRequest message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof proto.ResolverRequest
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {proto.ResolverRequest} ResolverRequest
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            ResolverRequest.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };
    
            /**
             * Verifies a ResolverRequest message.
             * @function verify
             * @memberof proto.ResolverRequest
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            ResolverRequest.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.id != null && message.hasOwnProperty("id"))
                    if (!$util.isString(message.id))
                        return "id: string expected";
                if (message.payload != null && message.hasOwnProperty("payload"))
                    if (!(message.payload && typeof message.payload.length === "number" || $util.isString(message.payload)))
                        return "payload: buffer expected";
                if (message.encrypted != null && message.hasOwnProperty("encrypted"))
                    if (typeof message.encrypted !== "boolean")
                        return "encrypted: boolean expected";
                if (message.publicKey != null && message.hasOwnProperty("publicKey"))
                    if (!(message.publicKey && typeof message.publicKey.length === "number" || $util.isString(message.publicKey)))
                        return "publicKey: buffer expected";
                return null;
            };
    
            /**
             * Creates a ResolverRequest message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof proto.ResolverRequest
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {proto.ResolverRequest} ResolverRequest
             */
            ResolverRequest.fromObject = function fromObject(object) {
                if (object instanceof $root.proto.ResolverRequest)
                    return object;
                var message = new $root.proto.ResolverRequest();
                if (object.id != null)
                    message.id = String(object.id);
                if (object.payload != null)
                    if (typeof object.payload === "string")
                        $util.base64.decode(object.payload, message.payload = $util.newBuffer($util.base64.length(object.payload)), 0);
                    else if (object.payload.length >= 0)
                        message.payload = object.payload;
                if (object.encrypted != null)
                    message.encrypted = Boolean(object.encrypted);
                if (object.publicKey != null)
                    if (typeof object.publicKey === "string")
                        $util.base64.decode(object.publicKey, message.publicKey = $util.newBuffer($util.base64.length(object.publicKey)), 0);
                    else if (object.publicKey.length >= 0)
                        message.publicKey = object.publicKey;
                return message;
            };
    
            /**
             * Creates a plain object from a ResolverRequest message. Also converts values to other types if specified.
             * @function toObject
             * @memberof proto.ResolverRequest
             * @static
             * @param {proto.ResolverRequest} message ResolverRequest
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            ResolverRequest.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.defaults) {
                    object.id = "";
                    if (options.bytes === String)
                        object.payload = "";
                    else {
                        object.payload = [];
                        if (options.bytes !== Array)
                            object.payload = $util.newBuffer(object.payload);
                    }
                    object.encrypted = false;
                    if (options.bytes === String)
                        object.publicKey = "";
                    else {
                        object.publicKey = [];
                        if (options.bytes !== Array)
                            object.publicKey = $util.newBuffer(object.publicKey);
                    }
                }
                if (message.id != null && message.hasOwnProperty("id"))
                    object.id = message.id;
                if (message.payload != null && message.hasOwnProperty("payload"))
                    object.payload = options.bytes === String ? $util.base64.encode(message.payload, 0, message.payload.length) : options.bytes === Array ? Array.prototype.slice.call(message.payload) : message.payload;
                if (message.encrypted != null && message.hasOwnProperty("encrypted"))
                    object.encrypted = message.encrypted;
                if (message.publicKey != null && message.hasOwnProperty("publicKey"))
                    object.publicKey = options.bytes === String ? $util.base64.encode(message.publicKey, 0, message.publicKey.length) : options.bytes === Array ? Array.prototype.slice.call(message.publicKey) : message.publicKey;
                return object;
            };
    
            /**
             * Converts this ResolverRequest to JSON.
             * @function toJSON
             * @memberof proto.ResolverRequest
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            ResolverRequest.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };
    
            /**
             * Gets the default type url for ResolverRequest
             * @function getTypeUrl
             * @memberof proto.ResolverRequest
             * @static
             * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns {string} The default type url
             */
            ResolverRequest.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
                if (typeUrlPrefix === undefined) {
                    typeUrlPrefix = "type.googleapis.com";
                }
                return typeUrlPrefix + "/proto.ResolverRequest";
            };
    
            return ResolverRequest;
        })();
    
        proto.ResolverResponse = (function() {
    
            /**
             * Properties of a ResolverResponse.
             * @memberof proto
             * @interface IResolverResponse
             * @property {string|null} [id] ResolverResponse id
             * @property {Uint8Array|null} [payload] ResolverResponse payload
             * @property {boolean|null} [encrypted] ResolverResponse encrypted
             */
    
            /**
             * Constructs a new ResolverResponse.
             * @memberof proto
             * @classdesc Represents a ResolverResponse.
             * @implements IResolverResponse
             * @constructor
             * @param {proto.IResolverResponse=} [properties] Properties to set
             */
            function ResolverResponse(properties) {
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }
    
            /**
             * ResolverResponse id.
             * @member {string} id
             * @memberof proto.ResolverResponse
             * @instance
             */
            ResolverResponse.prototype.id = "";
    
            /**
             * ResolverResponse payload.
             * @member {Uint8Array} payload
             * @memberof proto.ResolverResponse
             * @instance
             */
            ResolverResponse.prototype.payload = $util.newBuffer([]);
    
            /**
             * ResolverResponse encrypted.
             * @member {boolean} encrypted
             * @memberof proto.ResolverResponse
             * @instance
             */
            ResolverResponse.prototype.encrypted = false;
    
            /**
             * Creates a new ResolverResponse instance using the specified properties.
             * @function create
             * @memberof proto.ResolverResponse
             * @static
             * @param {proto.IResolverResponse=} [properties] Properties to set
             * @returns {proto.ResolverResponse} ResolverResponse instance
             */
            ResolverResponse.create = function create(properties) {
                return new ResolverResponse(properties);
            };
    
            /**
             * Encodes the specified ResolverResponse message. Does not implicitly {@link proto.ResolverResponse.verify|verify} messages.
             * @function encode
             * @memberof proto.ResolverResponse
             * @static
             * @param {proto.IResolverResponse} message ResolverResponse message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            ResolverResponse.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
                if (message.payload != null && Object.hasOwnProperty.call(message, "payload"))
                    writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.payload);
                if (message.encrypted != null && Object.hasOwnProperty.call(message, "encrypted"))
                    writer.uint32(/* id 3, wireType 0 =*/24).bool(message.encrypted);
                return writer;
            };
    
            /**
             * Encodes the specified ResolverResponse message, length delimited. Does not implicitly {@link proto.ResolverResponse.verify|verify} messages.
             * @function encodeDelimited
             * @memberof proto.ResolverResponse
             * @static
             * @param {proto.IResolverResponse} message ResolverResponse message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            ResolverResponse.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };
    
            /**
             * Decodes a ResolverResponse message from the specified reader or buffer.
             * @function decode
             * @memberof proto.ResolverResponse
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {proto.ResolverResponse} ResolverResponse
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            ResolverResponse.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.proto.ResolverResponse();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1: {
                            message.id = reader.string();
                            break;
                        }
                    case 2: {
                            message.payload = reader.bytes();
                            break;
                        }
                    case 3: {
                            message.encrypted = reader.bool();
                            break;
                        }
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };
    
            /**
             * Decodes a ResolverResponse message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof proto.ResolverResponse
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {proto.ResolverResponse} ResolverResponse
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            ResolverResponse.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };
    
            /**
             * Verifies a ResolverResponse message.
             * @function verify
             * @memberof proto.ResolverResponse
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            ResolverResponse.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.id != null && message.hasOwnProperty("id"))
                    if (!$util.isString(message.id))
                        return "id: string expected";
                if (message.payload != null && message.hasOwnProperty("payload"))
                    if (!(message.payload && typeof message.payload.length === "number" || $util.isString(message.payload)))
                        return "payload: buffer expected";
                if (message.encrypted != null && message.hasOwnProperty("encrypted"))
                    if (typeof message.encrypted !== "boolean")
                        return "encrypted: boolean expected";
                return null;
            };
    
            /**
             * Creates a ResolverResponse message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof proto.ResolverResponse
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {proto.ResolverResponse} ResolverResponse
             */
            ResolverResponse.fromObject = function fromObject(object) {
                if (object instanceof $root.proto.ResolverResponse)
                    return object;
                var message = new $root.proto.ResolverResponse();
                if (object.id != null)
                    message.id = String(object.id);
                if (object.payload != null)
                    if (typeof object.payload === "string")
                        $util.base64.decode(object.payload, message.payload = $util.newBuffer($util.base64.length(object.payload)), 0);
                    else if (object.payload.length >= 0)
                        message.payload = object.payload;
                if (object.encrypted != null)
                    message.encrypted = Boolean(object.encrypted);
                return message;
            };
    
            /**
             * Creates a plain object from a ResolverResponse message. Also converts values to other types if specified.
             * @function toObject
             * @memberof proto.ResolverResponse
             * @static
             * @param {proto.ResolverResponse} message ResolverResponse
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            ResolverResponse.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.defaults) {
                    object.id = "";
                    if (options.bytes === String)
                        object.payload = "";
                    else {
                        object.payload = [];
                        if (options.bytes !== Array)
                            object.payload = $util.newBuffer(object.payload);
                    }
                    object.encrypted = false;
                }
                if (message.id != null && message.hasOwnProperty("id"))
                    object.id = message.id;
                if (message.payload != null && message.hasOwnProperty("payload"))
                    object.payload = options.bytes === String ? $util.base64.encode(message.payload, 0, message.payload.length) : options.bytes === Array ? Array.prototype.slice.call(message.payload) : message.payload;
                if (message.encrypted != null && message.hasOwnProperty("encrypted"))
                    object.encrypted = message.encrypted;
                return object;
            };
    
            /**
             * Converts this ResolverResponse to JSON.
             * @function toJSON
             * @memberof proto.ResolverResponse
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            ResolverResponse.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };
    
            /**
             * Gets the default type url for ResolverResponse
             * @function getTypeUrl
             * @memberof proto.ResolverResponse
             * @static
             * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns {string} The default type url
             */
            ResolverResponse.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
                if (typeUrlPrefix === undefined) {
                    typeUrlPrefix = "type.googleapis.com";
                }
                return typeUrlPrefix + "/proto.ResolverResponse";
            };
    
            return ResolverResponse;
        })();
    
        proto.Execute = (function() {
    
            /**
             * Constructs a new Execute service.
             * @memberof proto
             * @classdesc Represents an Execute
             * @extends $protobuf.rpc.Service
             * @constructor
             * @param {$protobuf.RPCImpl} rpcImpl RPC implementation
             * @param {boolean} [requestDelimited=false] Whether requests are length-delimited
             * @param {boolean} [responseDelimited=false] Whether responses are length-delimited
             */
            function Execute(rpcImpl, requestDelimited, responseDelimited) {
                $protobuf.rpc.Service.call(this, rpcImpl, requestDelimited, responseDelimited);
            }
    
            (Execute.prototype = Object.create($protobuf.rpc.Service.prototype)).constructor = Execute;
    
            /**
             * Creates new Execute service using the specified rpc implementation.
             * @function create
             * @memberof proto.Execute
             * @static
             * @param {$protobuf.RPCImpl} rpcImpl RPC implementation
             * @param {boolean} [requestDelimited=false] Whether requests are length-delimited
             * @param {boolean} [responseDelimited=false] Whether responses are length-delimited
             * @returns {Execute} RPC service. Useful where requests and/or responses are streamed.
             */
            Execute.create = function create(rpcImpl, requestDelimited, responseDelimited) {
                return new this(rpcImpl, requestDelimited, responseDelimited);
            };
    
            /**
             * Callback as used by {@link proto.Execute#execute}.
             * @memberof proto.Execute
             * @typedef ExecuteCallback
             * @type {function}
             * @param {Error|null} error Error, if any
             * @param {proto.ResolverResponse} [response] ResolverResponse
             */
    
            /**
             * Calls Execute.
             * @function execute
             * @memberof proto.Execute
             * @instance
             * @param {proto.IResolverRequest} request ResolverRequest message or plain object
             * @param {proto.Execute.ExecuteCallback} callback Node-style callback called with the error, if any, and ResolverResponse
             * @returns {undefined}
             * @variation 1
             */
            Object.defineProperty(Execute.prototype.execute = function execute(request, callback) {
                return this.rpcCall(execute, $root.proto.ResolverRequest, $root.proto.ResolverResponse, request, callback);
            }, "name", { value: "Execute" });
    
            /**
             * Calls Execute.
             * @function execute
             * @memberof proto.Execute
             * @instance
             * @param {proto.IResolverRequest} request ResolverRequest message or plain object
             * @returns {Promise<proto.ResolverResponse>} Promise
             * @variation 2
             */
    
            return Execute;
        })();
    
        return proto;
    })();

    return $root;
})(protobuf);
