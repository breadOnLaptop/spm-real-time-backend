#ifndef API_CLIENT_H
#define API_CLIENT_H

#include <string>
#include <memory>
#include <grpcpp/grpcpp.h>
#include "telemetry.grpc.pb.h"

class ApiClient {
public:
    ApiClient(const std::string& endpoint);
    bool sendTelemetry(const telemetry::TelemetryPayload& payload);

private:
    std::string endpoint_;
};

#endif
