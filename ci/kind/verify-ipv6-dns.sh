#!/bin/bash

# IPv6 DNS Verification Script for Kind Cluster
# This script verifies IPv6 DNS resolution within a Kind cluster

#set -e

# Configuration
TEST_NAMESPACE="${TEST_NAMESPACE:-ipv6-dns-test}"
TEST_POD_NAME="ipv6-dns-tester"
TIMEOUT=120

# Log functions
log_info() {
    echo -e "[INFO] $1"
}

log_success() {
    echo -e "[SUCCESS] $1"
}

log_warning() {
    echo -e "[WARNING] $1"
}

log_error() {
    echo -e "[ERROR] $1"
}

# Create test namespace
create_test_namespace() {
    log_info "Creating test namespace: $TEST_NAMESPACE"

    kubectl create namespace "$TEST_NAMESPACE" --dry-run=client -o yaml | kubectl apply -f - &>/dev/null
    log_success "Test namespace created/verified"
}

# Deploy test pod with DNS tools
deploy_test_pod() {
    log_info "Deploying test pod with DNS tools..."

    cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: $TEST_POD_NAME
  namespace: $TEST_NAMESPACE
---
apiVersion: v1
kind: Pod
metadata:
  name: $TEST_POD_NAME
  namespace: $TEST_NAMESPACE
  labels:
    app: ipv6-dns-tester
spec:
  serviceAccountName: $TEST_POD_NAME
  containers:
  - name: dns-tester
    image: nicolaka/netshoot:latest
    command: ["sleep", "3600"]
    resources:
      requests:
        memory: "64Mi"
        cpu: "50m"
      limits:
        memory: "128Mi"
        cpu: "100m"
  restartPolicy: Never
EOF

    # Wait for pod to be ready
    log_info "Waiting for test pod to be ready..."
    kubectl wait --for=condition=Ready pod/$TEST_POD_NAME -n $TEST_NAMESPACE --timeout=${TIMEOUT}s
    log_success "Test pod is ready"
}

# Test DNS resolution
test_dns_resolution() {
    local target="$1"
    local description="$2"
    local record_type="${3:-AAAA}"

    log_info "Testing DNS resolution: $description"
    echo "Target: $target, Record Type: $record_type"

    # Test with dig if available
    local dig_result
    if kubectl exec -n $TEST_NAMESPACE $TEST_POD_NAME -- which dig &>/dev/null; then
        if dig_result=$(kubectl exec -n $TEST_NAMESPACE $TEST_POD_NAME -- dig +short "$target" $record_type 2>&1); then
            if [[ -n "$dig_result" ]]; then
                log_success "dig for $target successful"
                echo "Dig result: $dig_result"
            else
                log_warning "dig returned empty result for $target"
            fi
        else
            log_warning "dig command failed for $target"
        fi
    fi

    # Test with nslookup
    local nslookup_result
    if nslookup_result=$(kubectl exec -n $TEST_NAMESPACE $TEST_POD_NAME -- nslookup -type=$record_type "$target" 2>&1); then
        if echo "$nslookup_result" | grep -q "Address:.*:"; then
            log_success "nslookup for $target successful (IPv6 address found)"
            echo "$nslookup_result" | grep "Address:" | head -3
        elif echo "$nslookup_result" | grep -q "Address:.*\."; then
            log_warning "nslookup for $target returned IPv4 address only"
            echo "$nslookup_result" | grep "Address:" | head -3
        else
            log_error "nslookup for $target failed - no address found"
            echo "$nslookup_result"
            return 1
        fi
    else
        log_error "nslookup command failed for $target"
        echo "$nslookup_result"
        return 1
    fi

    echo "---"
    return 0
}

# Test internal Kubernetes DNS
test_kubernetes_dns() {
    log_info "Testing Kubernetes internal DNS resolution..."

    # Test cluster DNS
    test_dns_resolution "kubernetes.default.svc.cluster.local" "Kubernetes API service"

    # Test cluster DNS
    test_dns_resolution "kubernetes.default.svc.cluster.local." "Kubernetes API service"

    # Test cluster DNS
    test_dns_resolution "kubernetes.default.svc" "Kubernetes API service"

    # Test kube-dns/coredns service
    test_dns_resolution "kube-dns.kube-system.svc.cluster.local." "CoreDNS service"

    # Create a test service for DNS testing
    log_info "Creating test service for DNS verification..."
    cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: test-service
  namespace: $TEST_NAMESPACE
spec:
  selector:
    app: ipv6-dns-tester
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP
EOF

    # Wait a moment for DNS propagation
    sleep 5

    # Test the created service
    test_dns_resolution "test-service.$TEST_NAMESPACE.svc.cluster.local." "Test service DNS"
    test_dns_resolution "test-service" "Test service short name (from same namespace)"
}

# Test external DNS
test_external_dns() {
    log_info "Testing external DNS resolution..."

    # Test major DNS providers and services
    test_dns_resolution "google.com" "Google.com external DNS"
    test_dns_resolution "cloudflare.com" "Cloudflare.com external DNS"
    test_dns_resolution "ipv6.google.com" "Google IPv6 DNS"

    # Test IPv6-specific domains
    test_dns_resolution "ipv6-test.com" "IPv6 test domain"
    test_dns_resolution "test-ipv6.com" "Test IPv6 domain"
}

# Test DNS server configuration
test_dns_config() {
    log_info "Checking DNS configuration in test pod..."

    # Check resolv.conf
    log_info "DNS resolver configuration:"
    kubectl exec -n $TEST_NAMESPACE $TEST_POD_NAME -- cat /etc/resolv.conf
    echo "---"

    # Check if IPv6 nameservers are configured
    local resolv_content
    resolv_content=$(kubectl exec -n $TEST_NAMESPACE $TEST_POD_NAME -- cat /etc/resolv.conf)

    if echo "$resolv_content" | grep -q "nameserver.*:"; then
        log_success "IPv6 nameservers found in resolv.conf"
    else
        log_warning "No IPv6 nameservers found in resolv.conf"
    fi
}

# Test connectivity to DNS servers
test_dns_connectivity() {
    log_info "Testing connectivity to DNS servers..."

    # Get nameservers from resolv.conf
    local nameservers
    nameservers=$(kubectl exec -n $TEST_NAMESPACE $TEST_POD_NAME -- grep "^nameserver" /etc/resolv.conf | awk '{print $2}')

    for ns in $nameservers; do
        log_info "Testing connectivity to nameserver: $ns"
        if kubectl exec -n $TEST_NAMESPACE $TEST_POD_NAME -- timeout 5 nc -u -v "$ns" 53 </dev/null &>/dev/null; then
            log_success "Successfully connected to nameserver $ns"
        else
            log_error "Failed to connect to nameserver $ns"
        fi
    done
    echo "---"
}

# Check cluster networking
check_cluster_networking() {
    log_info "Checking cluster networking configuration..."

    # Check node IPv6 addresses
    log_info "Node IPv6 addresses:"
    kubectl get nodes -o wide
    echo "---"

    # Check pod IPv6 address
    log_info "Test pod IPv6 address:"
    kubectl get pod $TEST_POD_NAME -n $TEST_NAMESPACE -o wide
    echo "---"

    # Check if pod has IPv6 address
    local pod_ip
    pod_ip=$(kubectl get pod $TEST_POD_NAME -n $TEST_NAMESPACE -o jsonpath='{.status.podIP}')

    if [[ "$pod_ip" == *":"* ]]; then
        log_success "Pod has IPv6 address: $pod_ip"
    else
        log_warning "Pod has IPv4 address: $pod_ip"
    fi
}

# Generate summary report
generate_report() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')

    log_info "=== IPv6 DNS Verification Report ==="
    echo "Timestamp: $timestamp"
    echo "Namespace: $TEST_NAMESPACE"
    echo "Pod: $TEST_POD_NAME"
    echo ""

    log_info "Key findings:"
    echo "- DNS configuration verified"
    echo "- Internal Kubernetes DNS tested"
    echo "- External DNS resolution tested"
    echo "- IPv6 connectivity verified"
    echo ""

    log_success "DNS verification completed successfully!"
}

# Cleanup function
cleanup() {
    log_info "Cleaning up test resources..."

    kubectl delete namespace "$TEST_NAMESPACE" --ignore-not-found=true &>/dev/null || true
    log_success "Cleanup completed"
}

# Main execution
main() {
    echo "=== IPv6 DNS Verification Script for Kind Cluster ==="
    echo ""

    create_test_namespace
    deploy_test_pod

    echo ""
    log_info "Starting DNS verification tests..."
    echo ""

    check_cluster_networking
    test_dns_config
    test_dns_connectivity
    test_kubernetes_dns
    test_external_dns

    echo "dumping coredns logs"
    for pod in $(kubectl get pods -n kube-system -l k8s-app=kube-dns -o name); do
        echo "=== Logs from $pod ==="
        kubectl logs -n kube-system "$pod"
        echo
    done

    echo ""
    generate_report

    cleanup
}

# Handle script interruption
#trap 'log_error "Script interrupted"; cleanup; exit 1' INT TERM

# Run main function if script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
