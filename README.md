# Cloud Project : Bitly Clone

Bitly is a link management platform that lets you harness the power of your links by shortening, sharing, managing and analyzing links to your content. Billions of links are created every year by their millions of users, from individuals to small businesses to Fortune 500 companies. 
This project is a cloud scale Bitly-Like Service on Amazon Cloud with the web application deployed on Heroku.

![Screenshot from 2020-12-09 23-04-40](https://user-images.githubusercontent.com/14791021/101732992-f1106b80-3a72-11eb-874a-d4bcd6157cd2.png)

## Functionality
1. When a user goes to the URL provided by Heroku for web application, the user could enter the long URL and get the short URL on a button click.
2. Trending Data can be viewed i.e. the statistics about the links created by this application.
3. Number of hits and the last accessed time for each of the links can be viewed in the "Trending Data" page.

![Screenshot from 2020-12-09 23-17-10](https://user-images.githubusercontent.com/14791021/101734162-bad3eb80-3a74-11eb-9798-7eebde3c974c.png)

## AWS Services Used
1. EC2 Instances
2. VPC
3. Private/Public Subnets in VPC
4. Application Load Balancer
5. Network Load Balancer
6. Auto Scale Groups
7. Launch Configurations
8. Target Groups
9. Kong API Gateway

## Deployment 

1. Control Panel Servers (Private Subnet)  
These are dockerhost instances running a docker container with the Control Panel API which include the creating a short link API. These servers are created by auto scale groups which are linked with a launch configuration which has a custom ami of Control Panel and t2.micro instance type. These autoscale groups are then connected to target group (control-panel-target-group) which is linked to the internal application load balancer. The control-panel-target-group is associated by path ("/create*") i.e. any API call with the path starting with create will be redirected to Control Panel Server. 

2. Link Redirect Servers (Private Subnet)  
These are dockerhost instances running a docker container with the Link Redirect API which include the getting the long link API. These servers are created by auto scale groups which are linked with a launch configuration which has a custom ami of Link Redirect Server and t2.micro instance type. These autoscale groups are then connected to target group (link-redirect-target-group) which is linked to the internal application load balancer. The link-redirect-target-group is associated by path ("/get*") i.e. any API call with the path starting with get will be redirected to Control Panel Server.

3. NoSQL DB Cluster (Private Subnet)  
These are dockerhost instances running a docker container with NoSQL DB service which stores the trending data for each of the links. It consists of 5 nodes linked with each other using docker network "api_network" which is initialized as swarm and overlay.  All the nodes are then registered on each of the instances so that each one knows of the instances included in the cluster. With the help of this, if a document is created on Node 1, all the other nodes are notified and respective updates are made. These 5 nodes are connected to an internal network load balancer and only the Network Load Balancer's endpoint is used in communicating with the cluster.

4. MySQL Server (Private Subnet)  
This is a MySQL VM Instance that stores the short and the corresponding long links in a table. All the operations in this server are done by control-panel servers only.

5. RabbitMQ Server (Private Subnet)  
This is an EC2 Instance running the RabbitMQ service that acts as messaging queue service which stores all the messages coming from control panel server in a channel, in a queue and that queue is subscribed by the link redirect server and it consumes the messages and performs the appropriate tasks.

6. Kong API Gateway (Public Subnet)  
This is the public endpoint of the entire application for the web app. This adds an authentication layer to the application along with exposing it to the public. This is connected to the internal load balancer that connects the control panel APIs and Link Redirect APIs.


## Deployment Diagram
![Bitly Deployment Diagram](https://user-images.githubusercontent.com/14791021/101732759-92e38880-3a72-11eb-97e4-1430f270c06b.png)

## APIs

1. POST /create  

&nbsp;&nbsp; Request body:  
```json
{  
 "url" : "https://app.slack.com/client/T0AMW0A3S/C018XCQGGFL/thread/C018XCQGGFL-1607282036.001400"  
}  
```
&nbsp;&nbsp; Response :  
```json
&nbsp;&nbsp; {  
&nbsp;&nbsp; "long_URL": "https://app.slack.com/client/T0AMW0A3S/C018XCQGGFL/thread/C018XCQGGFL-1607282036.001400",  
&nbsp;&nbsp; "message": "201 Created",  
&nbsp;&nbsp; "short_URL": "http://cmpe.sjsu/osp8bx"  
&nbsp;&nbsp; }  
```

2. GET /getUrl
```json
&nbsp;&nbsp; Params:  
&nbsp;&nbsp; {  
&nbsp;&nbsp; "short_url" : "http://cmpe.sjsu/osp8bx"  
&nbsp;&nbsp; }
```

&nbsp;&nbsp; Response :
```json
&nbsp;&nbsp; {  
&nbsp;&nbsp; "hits": "2",  
&nbsp;&nbsp; "last_accessed": "2020-12-10 05:31:56.171214794 +0000 UTC m=+22922.251402975",    
&nbsp;&nbsp; "message": "301 Found",    
&nbsp;&nbsp; "url": "https://app.slack.com/client/T0AMW0A3S/C018XCQGGFL/thread/C018XCQGGFL-1607282036.001400"    
&nbsp;&nbsp; }  
```
   
