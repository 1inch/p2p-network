// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file relayer.proto (package relayer, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc, messageDesc } from "@bufbuild/protobuf/codegenv1";
import type { ResolverRequest, ResolverResponse } from "./resolver_pb";
import { file_resolver } from "./resolver_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file relayer.proto.
 */
export const file_relayer: GenFile = /*@__PURE__*/
  fileDesc("Cg1yZWxheWVyLnByb3RvEgdyZWxheWVyIjoKBUVycm9yEiAKBGNvZGUYASABKA4yEi5yZWxheWVyLkVycm9yQ29kZRIPCgdtZXNzYWdlGAIgASgJIlEKD0luY29taW5nTWVzc2FnZRISCgpwdWJsaWNLZXlzGAEgAygMEioKB3JlcXVlc3QYAiABKAsyGS5yZXNvbHZlci5SZXNvbHZlclJlcXVlc3QifwoPT3V0Z29pbmdNZXNzYWdlEi4KCHJlc3BvbnNlGAEgASgLMhoucmVzb2x2ZXIuUmVzb2x2ZXJSZXNwb25zZUgAEh8KBWVycm9yGAIgASgLMg4ucmVsYXllci5FcnJvckgAEhEKCXB1YmxpY0tleRgDIAEoDEIICgZyZXN1bHQqswEKCUVycm9yQ29kZRIeChpFUlJfSU5WQUxJRF9NRVNTQUdFX0ZPUk1BVBAAEh4KGkVSUl9SRVNPTFZFUl9MT09LVVBfRkFJTEVEEAESHQoZRVJSX0dSUENfRVhFQ1VUSU9OX0ZBSUxFRBACEiUKIUVSUl9SRVNQT05TRV9TRVJJQUxJWkFUSU9OX0ZBSUxFRBADEiAKHEVSUl9EQVRBX0NIQU5ORUxfU0VORF9GQUlMRUQQBEIsWipnaXRodWIuY29tLzFpbmNoL3AycC1uZXR3b3JrL3Byb3RvL3JlbGF5ZXJiBnByb3RvMw", [file_resolver]);

/**
 * Represents a standard error structure.
 *
 * @generated from message relayer.Error
 */
export type Error = Message<"relayer.Error"> & {
  /**
   * @generated from field: relayer.ErrorCode code = 1;
   */
  code: ErrorCode;

  /**
   * @generated from field: string message = 2;
   */
  message: string;
};

/**
 * Describes the message relayer.Error.
 * Use `create(ErrorSchema)` to create a new message.
 */
export const ErrorSchema: GenMessage<Error> = /*@__PURE__*/
  messageDesc(file_relayer, 0);

/**
 * IncomingMessage represents the message received via WebRTC data channel.
 *
 * @generated from message relayer.IncomingMessage
 */
export type IncomingMessage = Message<"relayer.IncomingMessage"> & {
  /**
   * @generated from field: repeated bytes publicKeys = 1;
   */
  publicKeys: Uint8Array[];

  /**
   * @generated from field: resolver.ResolverRequest request = 2;
   */
  request?: ResolverRequest;
};

/**
 * Describes the message relayer.IncomingMessage.
 * Use `create(IncomingMessageSchema)` to create a new message.
 */
export const IncomingMessageSchema: GenMessage<IncomingMessage> = /*@__PURE__*/
  messageDesc(file_relayer, 1);

/**
 * OutgoingMessage represents the response message to be sent via WebRTC data channel.
 *
 * @generated from message relayer.OutgoingMessage
 */
export type OutgoingMessage = Message<"relayer.OutgoingMessage"> & {
  /**
   * @generated from oneof relayer.OutgoingMessage.result
   */
  result: {
    /**
     * @generated from field: resolver.ResolverResponse response = 1;
     */
    value: ResolverResponse;
    case: "response";
  } | {
    /**
     * @generated from field: relayer.Error error = 2;
     */
    value: Error;
    case: "error";
  } | { case: undefined; value?: undefined };

  /**
   * @generated from field: bytes publicKey = 3;
   */
  publicKey: Uint8Array;
};

/**
 * Describes the message relayer.OutgoingMessage.
 * Use `create(OutgoingMessageSchema)` to create a new message.
 */
export const OutgoingMessageSchema: GenMessage<OutgoingMessage> = /*@__PURE__*/
  messageDesc(file_relayer, 2);

/**
 * Enum to represent standardized error codes.
 *
 * @generated from enum relayer.ErrorCode
 */
export enum ErrorCode {
  /**
   * Error in message format.
   *
   * @generated from enum value: ERR_INVALID_MESSAGE_FORMAT = 0;
   */
  ERR_INVALID_MESSAGE_FORMAT = 0,

  /**
   * Failed to resolve address for public key.
   *
   * @generated from enum value: ERR_RESOLVER_LOOKUP_FAILED = 1;
   */
  ERR_RESOLVER_LOOKUP_FAILED = 1,

  /**
   * gRPC execution failure.
   *
   * @generated from enum value: ERR_GRPC_EXECUTION_FAILED = 2;
   */
  ERR_GRPC_EXECUTION_FAILED = 2,

  /**
   * Failed to serialize the response.
   *
   * @generated from enum value: ERR_RESPONSE_SERIALIZATION_FAILED = 3;
   */
  ERR_RESPONSE_SERIALIZATION_FAILED = 3,

  /**
   * Failed to send the response via the data channel.
   *
   * @generated from enum value: ERR_DATA_CHANNEL_SEND_FAILED = 4;
   */
  ERR_DATA_CHANNEL_SEND_FAILED = 4,
}

/**
 * Describes the enum relayer.ErrorCode.
 */
export const ErrorCodeSchema: GenEnum<ErrorCode> = /*@__PURE__*/
  enumDesc(file_relayer, 0);

