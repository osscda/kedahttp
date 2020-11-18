#[macro_use]
extern crate serde_derive;

use kube::{
    api::{Object, Void, RawApi},
    client::APIClient,
    config,
};

use std::env;


#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct App {
    pub name: String,
    pub image: String,
    pub port: u32,
}

// This is a convenience alias that describes the object we get from Kubernetes
type KubeApp = Object<App, Void>;

fn main() {
    // Load the kubeconfig file.
    let kubeconfig = config::load_kube_config().expect("kubeconfig failed to load");

    // Create a new client
    let client = APIClient::new(kubeconfig);


    let namespace = std::env_var("KEDAHTTP_NAMESPACE").unwrap_or("default");

    // Describe the CRD we're working with.
    // This is basically the fields from our CRD definition.
    let resource = RawApi::customResource("App")
        .group("http.keda.sh")
        .within(&namespace);

    #[macro_use]
extern crate serde_derive;

use kube::{
    api::{Object, Void, RawApi},
    client::APIClient,
    config,
};


#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Book {
    pub title: String,
    pub authors: Option<Vec<String>>,
}

// This is a convenience alias that describes the object we get from Kubernetes
type KubeBook = Object<Book, Void>;

fn main() {
    // Load the kubeconfig file.
    let kubeconfig = config::load_kube_config().expect("kubeconfig failed to load");

    // Create a new client
    let client = APIClient::new(kubeconfig);

    // Set a namespace. We're just hard-coding for now.
    let namespace = "default";

    // Describe the CRD we're working with.
    // This is basically the fields from our CRD definition.
    let resource = RawApi::customResource("books")
        .group("example.technosophos.com")
        .within(&namespace);

    // Create our informer and start listening.
    let informer = Informer::raw(client, resource).init().expect("informer init failed");
    loop {
        informer.poll().expect("informer poll failed");

        // Now we just do something each time a new book event is triggered.
        while let Some(event) = informer.pop() {
            handle(event);
        }
    }
}

fn handle(event: WatchEvent<KubeApp>) {
    println!("Watch event for KubeApp:")
    dbg!(event)
}
