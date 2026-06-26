#include "api_client.h"
#include <iostream>

ApiClient::ApiClient(const std::string& endpoint) {
    endpoint_ = endpoint;
    // Strip port if present for the HTTP URL
    size_t colon_pos = endpoint_.find(":");
    if (colon_pos != std::string::npos) {
        endpoint_ = endpoint_.substr(0, colon_pos);
    }
}

bool ApiClient::sendTelemetry(const telemetry::TelemetryPayload& payload) {
    std::string binary_data = payload.SerializeAsString();
    
    std::string cmd = "curl -s -X POST -H 'Content-Type: application/octet-stream' --data-binary @- https://" + endpoint_ + "/api/telemetry/binary";
    
    FILE* pipe = popen(cmd.c_str(), "w");
    if (!pipe) {
        std::cerr << "Failed to open pipe for curl" << std::endl;
        return false;
    }
    
    fwrite(binary_data.data(), 1, binary_data.size(), pipe);
    int status = pclose(pipe);
    
    if (status == 0) {
        return true;
    } else {
        std::cerr << "HTTP POST failed with curl exit status: " << status << std::endl;
        return false;
    }
}
