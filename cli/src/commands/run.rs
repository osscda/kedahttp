use reqwest::{Client, Response, Error};
use std::collections::HashMap;
use std::result::Result;

pub async fn run(admin_url: &String, app_name: &String, image: &String, port: u32) 
-> Result<Response, Error> {

    let port_string = port.to_string();
    println!("run {} on port {}, named {}!", image, port, app_name);
    let client = Client::new();
    // TODO: maybe use Serde & a struct to serialize JSON
    let mut map = HashMap::new();
    map.insert("name", app_name);
    map.insert("image", &image);
    map.insert("port", &port_string);

    let request_url = format!("{}?name={}", admin_url, app_name);
    let res_future = client.post(&request_url)
    .json(&map)
    .send();

    res_future.await
}
