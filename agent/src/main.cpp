#include <iostream>
#include <thread>
#include <vector>
#include <chrono>
#include "api_client.h"
#include "telemetry_generator.h"

void agent_worker(int id, const std::string& endpoint) {
    std::string agent_id = "agent-node-00" + std::to_string(id);
    ApiClient client(endpoint);

    std::cout << "Starting " << agent_id << " connected to " << endpoint << std::endl;

    while (true) {
        telemetry::TelemetryPayload payload = TelemetryGenerator::generatePayload(agent_id);
        
        if (!client.sendTelemetry(payload)) {
            std::cerr << "[" << agent_id << "] Failed to send telemetry." << std::endl;
        } else {
            std::cout << "[" << agent_id << "] Telemetry sent successfully." << std::endl;
        }

        std::this_thread::sleep_for(std::chrono::seconds(2));
    }
}

int main(int argc, char* argv[]) {
    std::string endpoint = "spm_api:50051"; // Default gRPC endpoint
    if (argc > 1) {
        endpoint = argv[1];
    }

    std::vector<std::thread> workers;
    for (int i = 1; i <= 6; ++i) {
        workers.emplace_back(agent_worker, i, endpoint);
    }

    for (auto& w : workers) {
        w.join();
    }

    return 0;
}
