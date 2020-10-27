use structopt::StructOpt;
use std::fmt::Debug;

#[derive(Debug, StructOpt)]
enum Command {
    Rm {
        app_name: String,
    },
    Run {
        app_name: String,
        #[structopt(name="image", short)]
        image: String,
        #[structopt(name="port", short)]
        port: u32,
    },
}

#[derive(Debug, StructOpt)]
#[structopt(about = "Deploy scalable, production ready containers to Kubernetes")]
struct KedaHTTP {
    #[structopt(subcommand)]
    cmd: Command,
}

fn main() {
    let keda = KedaHTTP::from_args();
    
    match keda.cmd {
        Command::Rm{app_name} => {
            println!("remove {}!", app_name);
        },
        Command::Run{app_name, image, port} => {
            println!("run {} on port {}, named {}!", image, port, app_name);
        },
    }
}
