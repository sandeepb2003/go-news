package news

type Service struct {
    feedURL string
}

func NewService(feedURL string) *Service {
    return &Service{feedURL: feedURL}
}

func (s *Service) Headlines() []string {
    return []string{
        "No headlines loaded yet. Connect a data source.",
        "Project scaffold is ready for expansion.",
    }
}
