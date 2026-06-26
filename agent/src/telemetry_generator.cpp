#include "telemetry_generator.h"
#include <cstdlib>
#include <ctime>

std::map<std::string, int> TelemetryGenerator::agent_uptimes;

telemetry::TelemetryPayload TelemetryGenerator::generatePayload(const std::string& agent_id) {
    if (agent_uptimes.find(agent_id) == agent_uptimes.end()) {
        agent_uptimes[agent_id] = 3600 + (std::rand() % 86400);
    }
    agent_uptimes[agent_id] += 5;
    
    double cpu = (std::rand() % 1000) / 10.0;
    double mem = (std::rand() % 1000) / 10.0;
    double io = (std::rand() % 500) / 10.0;
    double ingress = (std::rand() % 1000) / 10.0;
    double egress = (std::rand() % 1000) / 10.0;
    
    double temp = 35.0 + (std::rand() % 450) / 10.0;
    std::string status = (temp > 70.0 || cpu > 90.0) ? "Warning" : "Healthy";
    int uptime = agent_uptimes[agent_id];

    telemetry::TelemetryPayload payload;
    payload.set_agent_id(agent_id);
    payload.set_cpu_utilization(cpu);
    payload.set_memory_utilization(mem);
    payload.set_disk_io(io);
    payload.set_network_ingress(ingress);
    payload.set_network_egress(egress);
    payload.set_temperature(temp);
    payload.set_uptime(uptime);
    payload.set_status(status);
    
    // Process 1
    auto* p1 = payload.add_top_processes();
    p1->set_pid(1);
    p1->set_executable_name("systemd");
    p1->set_resource_utilization(0.1);

    // Process 2
    auto* p2 = payload.add_top_processes();
    p2->set_pid(1000 + std::rand() % 9000);
    p2->set_executable_name("nginx");
    p2->set_resource_utilization((std::rand() % 150) / 10.0);

    // Process 3
    auto* p3 = payload.add_top_processes();
    p3->set_pid(1000 + std::rand() % 9000);
    p3->set_executable_name("postgres");
    p3->set_resource_utilization((std::rand() % 250) / 10.0);

    return payload;
}
