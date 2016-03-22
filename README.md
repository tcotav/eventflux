
## init your database

Use the cli `influx` to enter the following or whatever is similar for you:

    influx> CREATE DATABASE efdata
    influx> create user ef with password 'password'
    influx> grant all privileges to ef
    influx> CREATE RETENTION POLICY two_weeks ON efdata DURATION 2w REPLICATION 1 DEFAULT

You're set.  You've created a database, a user, and a retention policy.  Data is retained for up to the specificed duration.  See [the docs](https://docs.influxdata.com/influxdb/v0.10/query_language/database_management/#retention-policy-management) for more info. 


## set up your config

Copy the provided `etc/efconfig.json` to the appropriate location in your home directory.

    mkdir ~/.ef
    cp etc/efconfig.json ~/.ef
    vim ~/.ef/efconfig.json

change the values to match how you set up your database above.


## Storing an event

POST a json document of  EfEvents in the format below (mapping roughly to `influxdb` row format):

Json

    {
      "TagKey":"",
      "Tags":{
      },
      "Fields:{
      }
    }
   

Corresponds to the `golang` struct below:

    type EfEvent struct {
      TagKey string
      Tags   map[string]string
      Fields map[string]interface{}
    }


which in turn gets pushed into the influxdb database via a call in their v2 client:

    pt, err := client.NewPoint(event.TagKey, event.Tags, event.Fields, time.Now())

and then magically, we have stored our event in the database.


## Retrieving events
