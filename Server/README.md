# Set Up Your AWS CLI

## 1. **Install AWS CLI**:
If you haven't already, download and install the AWS CLI from the [AWS CLI installation page](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html).

## 2. **Configure the AWS CLI**:
```bash
aws configure
```

You'll be prompted to enter your AWS Access Key ID, Secret Access Key, default region, and default output format. Ensure you have the necessary permissions to create EKS clusters.

# Create an EKS Cluster

## 1. **Install eksctl**: 
`eksctl` is a simple tool to create and manage EKS clusters.

```bash
brew tap weaveworks/tap
brew install weaveworks/tap/eksctl
```

If you're not on macOS, follow the [installation guide for eksctl](https://eksctl.io/introduction/#installation).

## 2. **Create a new EKS Cluster**: 
```bash
eksctl create cluster --name <cluster-name> --version 1.25 --region <region> --nodegroup-name linux-nodes --node-type t3.medium --nodes 3 --nodes-min 1 --nodes-max 4 --managed
```

This command will take some time to complete, as it sets up all the necessary resources in AWS for your EKS cluster.


## 3. Update kubeconfig for kubectl

Once your EKS cluster is created, you'll need to update your `kubeconfig` so that `kubectl` can communicate with the new cluster.

1. **Update kubeconfig**:
   ```bash
   aws eks --region ca-central-1 update-kubeconfig --name <cluster-name>
   ```

2. **Verify the configuration**:
   ```bash
   kubectl get nodes
   ```
   You should see a list of nodes indicating that `kubectl` is correctly set to work with your EKS cluster.


# Deploy Temporal

With your EKS cluster ready, you can now deploy Temporal. For learning purposes, you can use the provided Temporal Helm chart.

1. **Install Helm**: 
Ensure Helm is installed on your machine. If not, you can install it using:

```bash
brew install helm
```

2. **Deploy Temporal using Helm**:
Add the Temporal Helm repository and deploy Temporal:

```bash
helm repo add temporalio https://temporalio.github.io/helm-charts
helm install temporal-training temporalio/temporal --version 0.52.0
```

3. **Verify Installation**:
Check the Temporal services within your cluster:

```bash
kubectl get pods
kubectl get services
```

# Troubleshooting
## Unable to connect to the cluster
### IAM and cli access

Ensure that the IAM role used has the necessary permissions to interact with the EKS cluster.

- **IAM Permissions**: Check that your IAM role has `eks:DescribeCluster`, `eks:ListClusters`, and other necessary permissions.
- **Authentication and Authorization**: Make sure that your local AWS CLI is properly authenticated to access AWS services (`aws configure`).

#### Network and Security Groups

EKS and Lens require proper networking configurations which include VPC, subnets, and security groups.

- **Security Groups**: Ensure that the security groups associated with your cluster nodes allow inbound and outbound traffic on necessary ports (e.g., 443 for HTTPS).
- **Network ACLs**: Check that any network ACLs do not block the traffic needed to/from the cluster.
### Cluster and Node Group Status

## 'default' namespace not found
### 1. Install tctl
```bash
brew install tctl
```
### 2. Port forward the 7233 frontend pod
```bash
kubectl port-forward pod/<frontend pod name> 7233 7233
```
### 3. Create the default namespace
```bash
tctl --ns default namespace register -rd 3
```

Check the status of your EKS cluster and the node groups to ensure they are active and healthy.

- **Cluster Status**:
  ```bash
  aws eks --region us-west-2 describe-cluster --name <cluster-name> --query cluster.status
  ```
- **Node Group Status**:
  ```bash
  eksctl get nodegroup --cluster <cluster-name> --region us-west-2
  ```
