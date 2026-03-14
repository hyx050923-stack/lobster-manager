func NewManager() (*Manager, error) {

	switch runtime.GOOS {

	case "android":

		return &Manager{
			runtime: runtime.NewUDocker(),
		}, nil

	default:

		cli, err := client.NewClientWithOpts(
			client.FromEnv,
			client.WithAPIVersionNegotiation(),
		)

		if err != nil {
			return nil, err
		}

		return &Manager{
			cli: cli,
			ctx: context.Background(),
		}, nil
	}
}