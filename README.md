# Requirements for the task

Implement URL shortener web-application, using programming language, databases and frameworks of your chioce.
Index page should display a form offering user to shorten a URL. After form was submitted, short URL should be displayed to the user.

Short URLs should redirect to original URLs.
Your implementation should be able to sustain 100 URL submissions and 1,000 short URL accesses per minute.

You may assume that there will be approximately 1,000,000 unique URLs and all further submissions will be duplicates of URLs submitted previously.

# Shorten-url
Simple Shorten Url

## Requirements
1. Docker
2. Docker-compose
# Installation
1. Clone project git clone https://github.com/shevchik87/shorten-url.git
2. cd shorten-url && run docker-compose up (server will start on http://127.0.0.1:80)
3. enjoy 😊 


# Scale

What we have to do that scale up our app?

1) First of all we add cache (Redis) to our project. In redis we will write hot URLs that are frequently accessed. If we follow the 80-20 rule meaning 20% of URLs generate 80% of traffic. We need store in cache only 20% URLs. Of course we need a cluster of Redis, which will be able to handle such a load.

2) To scale out our DB, we need to partition it so that it can store information about billions of URLs. We need to come up with a partitioning scheme that would divide and store our data to different DB servers. For this we need apply Hash-Based Partitioning: In this scheme, we take a hash of the object we are storing. We then calculate which partition to use based upon the hash. In our case, we can take the hash of the ‘key’ or the actual URL to determine the partition in which we store the data object.Our hashing function will randomly distribute URLs into different partitions (e.g., our hashing function can always map any key to a number between [1…256]), and this number would represent the partition in which we store our object.

3) Of course we need to setup load balancer between Clients and Application servers



