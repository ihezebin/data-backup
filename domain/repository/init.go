package repository

func Init() {
	SetExampleRepository(NewExampleMongoRepository())
}
