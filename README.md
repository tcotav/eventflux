
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
