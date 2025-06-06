// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file resolver.proto (package resolver, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file resolver.proto.
 */
export const file_resolver: GenFile = /*@__PURE__*/
  fileDesc("Cg5yZXNvbHZlci5wcm90bxIIcmVzb2x2ZXIiOwoFRXJyb3ISIQoEY29kZRgBIAEoDjITLnJlc29sdmVyLkVycm9yQ29kZRIPCgdtZXNzYWdlGAIgASgJIlQKD1Jlc29sdmVyUmVxdWVzdBIKCgJpZBgBIAEoCRIRCgllbmNyeXB0ZWQYAiABKAgSDwoHcGF5bG9hZBgDIAEoDBIRCglwdWJsaWNLZXkYBCABKAwicAoQUmVzb2x2ZXJSZXNwb25zZRIKCgJpZBgBIAEoCRIRCgllbmNyeXB0ZWQYAiABKAgSEQoHcGF5bG9hZBgDIAEoDEgAEiAKBWVycm9yGAQgASgLMg8ucmVzb2x2ZXIuRXJyb3JIAEIICgZyZXN1bHQqbgoJRXJyb3JDb2RlEhoKFkVSUl9JTlRFUk5BTF9FWENFUFRJT04QABIeChpFUlJfSU5WQUxJRF9NRVNTQUdFX0ZPUk1BVBABEiUKIUVSUl9SRVNQT05TRV9TRVJJQUxJWkFUSU9OX0ZBSUxFRBACMksKB0V4ZWN1dGUSQAoHRXhlY3V0ZRIZLnJlc29sdmVyLlJlc29sdmVyUmVxdWVzdBoaLnJlc29sdmVyLlJlc29sdmVyUmVzcG9uc2VCLVorZ2l0aHViLmNvbS8xaW5jaC9wMnAtbmV0d29yay9wcm90by9yZXNvbHZlcmIGcHJvdG8z");

/**
 * Represents a standard error structure.
 *
 * @generated from message resolver.Error
 */
export type Error = Message<"resolver.Error"> & {
  /**
   * @generated from field: resolver.ErrorCode code = 1;
   */
  code: ErrorCode;

  /**
   * @generated from field: string message = 2;
   */
  message: string;
};

/**
 * Describes the message resolver.Error.
 * Use `create(ErrorSchema)` to create a new message.
 */
export const ErrorSchema: GenMessage<Error> = /*@__PURE__*/
  messageDesc(file_resolver, 0);

/**
 * @generated from message resolver.ResolverRequest
 */
export type ResolverRequest = Message<"resolver.ResolverRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: bool encrypted = 2;
   */
  encrypted: boolean;

  /**
   * @generated from field: bytes payload = 3;
   */
  payload: Uint8Array;

  /**
   * @generated from field: bytes publicKey = 4;
   */
  publicKey: Uint8Array;
};

/**
 * Describes the message resolver.ResolverRequest.
 * Use `create(ResolverRequestSchema)` to create a new message.
 */
export const ResolverRequestSchema: GenMessage<ResolverRequest> = /*@__PURE__*/
  messageDesc(file_resolver, 1);

/**
 * @generated from message resolver.ResolverResponse
 */
export type ResolverResponse = Message<"resolver.ResolverResponse"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: bool encrypted = 2;
   */
  encrypted: boolean;

  /**
   * @generated from oneof resolver.ResolverResponse.result
   */
  result: {
    /**
     * @generated from field: bytes payload = 3;
     */
    value: Uint8Array;
    case: "payload";
  } | {
    /**
     * @generated from field: resolver.Error error = 4;
     */
    value: Error;
    case: "error";
  } | { case: undefined; value?: undefined };
};

/**
 * Describes the message resolver.ResolverResponse.
 * Use `create(ResolverResponseSchema)` to create a new message.
 */
export const ResolverResponseSchema: GenMessage<ResolverResponse> = /*@__PURE__*/
  messageDesc(file_resolver, 2);

/**
 * Enum to represent standardized error codes.
 *
 * @generated from enum resolver.ErrorCode
 */
export enum ErrorCode {
  /**
   * gRPC execution failure.
   *
   * @generated from enum value: ERR_INTERNAL_EXCEPTION = 0;
   */
  ERR_INTERNAL_EXCEPTION = 0,

  /**
   * Error in message format.
   *
   * @generated from enum value: ERR_INVALID_MESSAGE_FORMAT = 1;
   */
  ERR_INVALID_MESSAGE_FORMAT = 1,

  /**
   * Failed to serialize the response.
   *
   * @generated from enum value: ERR_RESPONSE_SERIALIZATION_FAILED = 2;
   */
  ERR_RESPONSE_SERIALIZATION_FAILED = 2,
}

/**
 * Describes the enum resolver.ErrorCode.
 */
export const ErrorCodeSchema: GenEnum<ErrorCode> = /*@__PURE__*/
  enumDesc(file_resolver, 0);

/**
 * @generated from service resolver.Execute
 */
export const Execute: GenService<{
  /**
   * @generated from rpc resolver.Execute.Execute
   */
  execute: {
    methodKind: "unary";
    input: typeof ResolverRequestSchema;
    output: typeof ResolverResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_resolver, 0);

