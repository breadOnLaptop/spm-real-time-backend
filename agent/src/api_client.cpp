#include "api_client.h"
#include <iostream>

ApiClient::ApiClient(const std::string& endpoint) {
    auto channel = grpc::CreateChannel(endpoint, grpc::InsecureChannelCredentials());
    stub_ = telemetry::TelemetryService::NewStub(channel);
}

bool ApiClient::sendTelemetry(const telemetry::TelemetryPayload& payload) {
    telemetry::TelemetryResponse response;
    grpc::ClientContext context;
    
    // timeout for the rpc call
    gpr_timespec deadline;
    deadline.tv_sec = 2;
    deadline.tv_nsec = 0;
    deadline.clock_type = GPR_TIMESPAN;
    context.set_deadline(deadline);

    grpc::Status status = stub_->SendTelemetry(&context, payload, &response);
    
    if (status.ok()) {
        return response.success();
    } else {
        std::cerr << "gRPC failed: " << status.error_code() << ": " << status.error_message() << std::endl;
        return false;
    }
}
