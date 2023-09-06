// @generated by protoc-gen-connect-es v0.13.2 with parameter "target=ts"
// @generated from file wg/cosmo/node/v1/node.proto (package wg.cosmo.node.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { GetConfigRequest, GetConfigResponse } from "./node_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service wg.cosmo.node.v1.NodeService
 */
export const NodeService = {
  typeName: "wg.cosmo.node.v1.NodeService",
  methods: {
    /**
     * @generated from rpc wg.cosmo.node.v1.NodeService.GetLatestValidRouterConfig
     */
    getLatestValidRouterConfig: {
      name: "GetLatestValidRouterConfig",
      I: GetConfigRequest,
      O: GetConfigResponse,
      kind: MethodKind.Unary,
    },
  }
} as const;

