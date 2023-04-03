+++
#title = "Parigot"
description = "Where you go to anything related to parigot, the operating environment for WASM microservices."
weight = 2
disableReadmoreNav = true
toc = false
+++


<div class="row w-80">
    <blockquote class="display-4">parigot: the operating environment for microservices</blockquote>
</div>
<div class="row justify-content-center align-items-center w-80">
    <div class="col text-left border-bottom-0">
    <h4>For developers</h4>
    <p>
    parigot is both an operating environment and a development tool.  parigot offers you the ability to develop, test,
    and debug your microservice apps as a monolith and then deploy them as (normal) microservices.  You can probably
    see what this means:
    <ul>
        <li>Reproducible test results</li>
        <li>Well-defined startup ordering</lib>
        <li>Intra-service dependency resolution (including detecting cycles)</lib>
        <li>Debuggers work because in development your "app" is a single process</lib>
        <li>Unified logging, not one log per service</li>
    </ul>
    </p>
    <h4>How does it work?</h4>
    <p>
    Fundamentally, parigot works because it has a different programming model.  The programming model allows parigot's 
    infrastructure to make significant assumptions about how the program works.  For example, there is no visiblity or
    access to the "raw" network... parigot's model only allows calls to other services.  (Unless you build your own services that
    offer this capability!) Thus, the infrastructure can make assumptions about how the program works (e.g. listening
    on sockets) and can then take different actions than a typical development model.  A typical development model
    would mean that the infrastructure cannot make assumptions about what ports are in use, which ones are connections
    to other services in the app, etc.
    </p>
    <h4>
    What languages do you support writing parigot programs in?
    </h4>
    <p>
    In principle, all of them.  parigot's infrastructure is built using <a ref="https://webassembly.org">Web Assembly</a> (or 
    <a ref="https://docs.docker.com/desktop/wasm/">wasm-based containers</a>) so
    <a ref="https://github.com/appcypher/awesome-wasm-langs"> any language that compile to WASM </a> will work.  
    </p>
    <p>
    Each language,
    however, requires some glue code to hook up to parigot's ABI; this is exactly analagous to various languages that need to either use
    or replicate the C language library libc.a for linux.  This library, in both the parigot and linux case, effectively
    sets up requests for the kernel by laying out some bytes in memory and then requesting a kernel service.  The kernel
    looks at the bytes (the parameters), performs the action, and then writes some bytes back into memory that are the
    result.  We have support for golang already working and python is next.
    </p>
    <h4>How do we make money</h4>
    <p>
    When the software is open source, of course folks ask this question a great deal.  We are developing a deploy and hosting
    service that will run parigot-based programs at a significantly lower price than any other cloud provider. So, a parigot
    program can be developed on your local laptop in the "easy way" as a single process and then deployed into to the
    cloud as discrete services that talk to each other over the network.  
    </p>
    <p>
    You can still parigot yourself, with your own programs and services, without deploying to our service.
    </p>
        <h4>
    What does programming model actually look like?
    </h4>
    <p>
    Here's a snippet of a simple program, written in go:
    </p>
    <p class="white-space:pre; font-family: "Courier New", monospace;">

func main() {

    lib.FlagParseCreateEnv()

    if _, err := callImpl.Require1("log", "LogService"); err != nil {
        panic("unable to require log service: " + err.Error())
    }
	if _, err := callImpl.Require1("file", "FileService"); err != nil {
		panic("unable to require file service:" + err.Error())
	}
    // This code is the *implementation* of the StoreService.
	if _, err := callImpl.Export1("demo.vvv", "StoreService"); err != nil {
		panic("unable to export demo.vvv: " + err.Error())
	}
	store.RunStoreService(&myStoreServer{})
    //...
}

    // Ready is a check, if this returns false the library should abort and not attempt to run this service.
    // Normally, this is used to block using the lib.Run() call.  This call will wait until all the required
    // services are ready.


func (m *myStoreServer) Ready() bool {

	if _, err := callImpl.Run(&pbsys.RunRequest{Wait: true}); err != nil {
		print("ready: error in attempt to signal Run: ", err.Error(), "\n")
		return false
	}

	logger, err := log.LocateLogService()
	if err != nil {
		print("ERROR trying to create log client: ", err.Error(), "\n")
		return false
	}
	m.logger = logger

	fClient, err := file.LocateFileService(logger)
	if err != nil {
		print("ERROR trying to create fs client: ", err.Error(), "\n")
		return false
	}
    // ..
}
</p>
    <h5>Discussion</h5>
    <p>Some things to note from this example:
    <ul>
    <li>You have explicitly require and export services you depend on or provide, respectively.
    <li> parigot can build a model of everything that is needed to run your program--and maybe more importantly what is not
        needed.  If your system has 300 services in production, it is unlikely you need all 300 to run your program and
        parigot can compute what is needed exactly.  </li>
    <li>Secondly, this give parigot a dependency ordering.  In other words,
        we make sure that your code will not start until everything it requires is up and ready.  If there is no 
        deterministic startup order, parigot will abort.</li>
    <li>The function `Ready()` is called to inform your program that your dependencies are prepared.  You then can
    call LocateSOMETHING() to get an interface that is the implementation of that service.  Note that you cannot
    tell if that service is in the same address space or over the network! </li>
    </ul>
    </p>
    </div>
</div>

<div class="row w-80">
    <p class="blockquote-footer display-4"><strong>did we forget to mention it is free and open source?</strong></p>
</div>
<script lang="javascript">
    console.log("since you went to the trouble to look here...");
    console.log("we give our releases names based on the letters of the alphabet, with each release being a city.");
    console.log("the first release... not the 1.0, just the first public release is called 'atlanta'.");
    console.log("second release name not totally decided yet, but possibly 'boradino' or 'barcalona'.");
</script>