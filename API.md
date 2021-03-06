

# API

The backend exposes an API at `http://localhost:8080/api`.

<details>
<summary>List resources</summary>

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/resources](http://localhost:8080/api/resources)  | POST  | Return list of cloud resources |  :white_check_mark: |


To filter the resources, send a body containing a query.

```js
{
  // set a limit, default is 25, max is 2000 
  "limit": 25,
  //specifies the number of rows to skip before any rows are retrieved
  "offset": 0,
  //filter the resources
  "filter": {
    "core.type": "ec2.Instance"
    "tags.env": "prod"
  }
  //optional sort
  "sort": ["core.type"]
}
```

The response contains:

```js
{
  //the total number of resources matching the filter
  "count": 25,
  //the fields available for filtering the resources, their count is updated based on the input query
  //a field can be: a resource type, a region, a tag key
  "fieldGroups": [],
  //the resources matching the input query (paginated)
  "resources": []
}
```

- Examples of requests:
```js
{
  // set a limit, default is 25, max is 100 
  "limit": 25,
  //specifies the number of rows to skip before any rows are retrieved
  "offset": 0,
  //filter the resources
  "filter": {
    "core.type": "ec2.Instance"
  }
  //optional sort
  "sort": ["core.type"]
}
```

```shell
curl --location --request POST 'http://localhost:8080/api/resources' \
--header 'Content-Type: application/json' \
--data-raw '{
    "filter": {
        "tags.kubernetes.io/created-for/pv/name": "opta-persistent-0-hellopv-hellopv-k8s-service-0"
    }
}'
```

Example of queries:
```js

//default return all the resources (no payload)
{}

//return resources of type "ec2.Instance" with the tag "team" equals "marketplace"
{
  "filter":{
    "core.type": "ec2.Instance",
	  "tags.team": "marketplace"
  }
}

//return resources with the tag "team" defined
{
  "filter":{
	  "tags.team": "(not null)"
  }
}

//return resources missing the tag "team"
{
  "filter":{
	  "tags.team": "(missing)"
  }
}

//filter with more than one value for a field using a OR
// will return resources with type=ec2.Volume AND (team="marketplace" OR team="shipping")
{
  "filter":{
    "core.type":"ec2.Volume",
    "$or": [
      { "tags.team": "marketplace" },
      { "tags.team": "shipping" }
    ]
  } 
}

//Using multiple OR sections
// will return resources with (team="marketplace" OR team="shipping") AND (cluster="dev" OR cluster="prod")  AND (size="large" OR size="medium") 
{
  "filter":{
    "$or": [
      { "tags.team": "marketplace" },
      { "tags.team": "shipping" }
    ],
    "$and": [
      { "$or": [
        { "tags.cluster": "dev" },
        { "tags.cluster": "prod" }
      ] },
      { "$or": [
        { "tags.size": "large" },
        { "tags.size": "medium" }
      ] }
    ]
  }
}

//sort by a field
{
  "filter":{
    "core.type": "s3.Bucket"
  },
  "sort": ["core.region"]
}

//The default order for column is ascending order but you can control it with an optional prefix: + or -. + means ascending order, and - means descending order.
//sort by region desc
{
  "filter":{
    "core.type": "s3.Bucket"
  },
  "sort": ["-core.region"]
}

//Set a limit: default 25, Max is 100
//return the ec2.Instance with a limit of 10 results
{
  "limit": 10,
  "filter":{
    "core.type": "ec2.Instance"
  }
}

//used with limit, the offset paramerter specifies the number of rows to skip before any rows are retrieved
//first page: first 10 results
{
  "limit": 10,
  "offset": 0,
  "filter":{
    "core.type": "ec2.Instance"
  }
}
//second page: next 10 results
{
  "limit": 10,
  "offset": 10,
  "filter":{
    "core.type": "ec2.Instance"
  }
}

```

- Examples of response:

```js
{
  //the query has 2 results - if pagination is used this number would be unchanged - it's the total number of results
  "count":2,
  //the updated list of fields matching the current query
  "fieldGroups":[
    {
      "name":"core",
      "fields":[
        {
          "name":"region",
          "values":[
            {
              "value":"us-east-1",
              "count":"2"
            },
            {
              "value":"eu-west-3",
              // "-" means that this value is not relevant for current query - filtering on it would have no effect
              "count":"-"
            }
          ],
          "count":2
        },
        {
          "name":"type",
          "values":[
            {
              "value":"ec2.Instance",
              "count":"2"
            },
            {
              "value":"ec2.Volume",
              "count":"-"
            }
          ],
          "count":2
        }
      ]
    },
    {
      "name":"tags",
      "fields":[
        {
          "name":"managed-by",
          "values":[
            {
              "value":"cloudformation",
              "count":"2"
            },
            {
              "value":"terraform",
              "count":"-"
            }
          ],
          "count":2
        },
        {
          "name":"env",
          "values":[
            {
              "value":"prod",
              "count":"1"
            },
            {
              "value":"dev",
              "count":"1"
            }
          ],
          "count":2
        }
      ]
    }
  ],
  //list the resources - if pagination is used only one page would be returned at time
  "resources":[
    {
      "id":"i-05a8cc7c8b7bc4f2d",
      "region":"us-east-1",
      "type":"ec2.Instance",
      "tags":[
        {
          "key":"env",
          "value":"dev"
        },
        {
          "key":"managed-by",
          "value":"cloudformation"
        }
      ],
      "rawData":{
        "AmiLaunchIndex":0,
        "Architecture":"x86_64",
        "BlockDeviceMappings":[
          {
            "DeviceName":"/dev/xvda",
            "Ebs":{
              "AttachTime":"2022-06-16T23:25:00Z",
              "DeleteOnTermination":true,
              "Status":"attached",
              "VolumeId":"vol-0d125183ed4159484"
            }
          }
        ],
        "BootMode":"",
        "ImageId":"ami-032930428bf1abbff",
        "InstanceId":"i-05a8cc7c8b7bc4f2d"
      }
    },
    {
      "id":"i-0695984d3a9256cea",
      "region":"us-east-1",
      "type":"ec2.Instance",
      "tags":[
        {
          "key":"env",
          "value":"dev"
        },
        {
          "key":"managed-by",
          "value":"cloudformation"
        }
      ],
      "rawData":{
        "AmiLaunchIndex":0,
        "Architecture":"x86_64",
        "BlockDeviceMappings":[
          {
            "DeviceName":"/dev/xvda",
            "Ebs":{
              "AttachTime":"2022-06-16T23:25:00Z",
              "DeleteOnTermination":true,
              "Status":"attached",
              "VolumeId":"vol-0f6ee55f46d5b5f65"
            }
          }
        ],
        "BootMode":"",
        "ImageId":"ami-032930428bf1abbff",
        "InstanceId":"i-0695984d3a9256cea"
      },
      "updatedAt":"2022-06-20T14:10:44.679424-07:00"
    }
  ]
}
```


</details>
<details>
<summary>Get a resource</summary>

| Route | Method |  Description |  Status |
| ------------- | ------------- | ------------- | ------------- |
| [/resource](http://localhost:8080/api/resource)  | GET  | Return a resource |  :white_check_mark: |

| Parameters | Description |  Examples |
| ------------- | ------------- | ------------- |
| id  | the resource id  | `id=i-024c4971f7f510c8f` return resource with the id `i-024c4971f7f510c8f`

</details>
<details>
<summary>Get Engine Status</summary>

Returns the Status of the Cloudgrep run.

| Route                                                   | Method | Description              |  Status |
|---------------------------------------------------------| ------------- |--------------------------| ------------- |
| [/enginestatus](http://localhost:8080/api/enginestatus) | GET  | Return the Engine status |  :white_check_mark: |

Sample Responses:
```js
// Engine completed successfully
{
    "runId": "6fd67489-d852-4962-95bc-eea01159993f",
    "eventType": "engine",
    "status": "success",
    "providerName": "",
    "resourceType": "",
    "error": "",
    "createdAt": "2022-06-22T02:54:12.727066+05:30",
    "updatedAt": "2022-06-22T02:54:25.458235+05:30",
    "childEvents": [
        {
            "runId": "6fd67489-d852-4962-95bc-eea01159993f",
            "eventType": "provider",
            "status": "success",
            "providerName": "aws",
            "resourceType": "",
            "error": "",
            "createdAt": "2022-06-22T02:54:12.727395+05:30",
            "updatedAt": "2022-06-22T02:54:13.979699+05:30",
            "childEvents": null
        },
        {
            "runId": "6fd67489-d852-4962-95bc-eea01159993f",
            "eventType": "resource",
            "status": "success",
            "providerName": "AWS Provider for account 693658092572, region us-east-2",
            "resourceType": "ec2.Volume",
            "error": "",
            "createdAt": "2022-06-22T02:54:13.980207+05:30",
            "updatedAt": "2022-06-22T02:54:16.658743+05:30",
            "childEvents": null
        }
    ]
}

// Engine is currently running
{
    "runId": "6fd67489-d852-4962-95bc-eea01159993f",
    "eventType": "engine",
    "status": "failed",
    "providerName": "",
    "resourceType": "",
    "error": "1 error message\n error message",
    "createdAt": "2022-06-22T02:54:12.727066+05:30",
    "updatedAt": "2022-06-22T02:54:25.458235+05:30",
    "childEvents": [
    {
        "runId": "6fd67489-d852-4962-95bc-eea01159993f",
        "eventType": "provider",
        "status": "success",
        "providerName": "aws",
        "resourceType": "",
        "error": "",
        "createdAt": "2022-06-22T02:54:12.727395+05:30",
        "updatedAt": "2022-06-22T02:54:13.979699+05:30",
        "childEvents": null
    },
    {
        "runId": "6fd67489-d852-4962-95bc-eea01159993f",
        "eventType": "resource",
        "status": "failed",
        "providerName": "AWS Provider for account 693658092572, region us-east-2",
        "resourceType": "ec2.Volume",
        "error": "error message",
        "createdAt": "2022-06-22T02:54:13.980207+05:30",
        "updatedAt": "2022-06-22T02:54:16.658743+05:30",
        "childEvents": null
    }
]
}
```

If you need to know when the engine is done running, keep pulling this endpoint until the status is no longer **fetching**.

</details>
<details>
<summary>Refresh the resources</summary>

Trigger the engine to refresh the cloud resources.
Calling this endpoint will returns immediately, the engine will start fetching the resources async.

| Route                                                   | Method | Description              |  Status |
|---------------------------------------------------------| ------------- |--------------------------| ------------- |
| [/refresh](http://localhost:8080/api/refresh) | POST  | Refresh the cloud resources |  :white_check_mark: |

Sample Responses:
```js
// Refresh request acknowledged, the refresh has started.
code: 200
body: {}

// The refresh has already been triggered and is in progress
code: 202
{
  "status":"202",
  "error":"engine is already running"
}

// There was an error
code: 400
{
  "status":"400",
  "error":"can't connect to datastore"
}

```

Once the refreshed is triggered, call **Get Engine Status** API to know if the refresh is done.
</details>
