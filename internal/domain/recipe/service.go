package recipe

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Service struct {
	client *http.Client
}

func NewService() *Service {
	return &Service{
		client: &http.Client{},
	}
}

func (s *Service) GetRecipe(ctx context.Context, url string) (*Recipe, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recipe: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	r := &Recipe{}

	// Extract recipe data using Schema.org markup
	doc.Find("[itemtype='http://schema.org/Recipe']").Each(func(i int, s *goquery.Selection) {
		r.Title = s.Find("[itemprop='name']").Text()
		r.CookTime = s.Find("[itemprop='cookTime']").Text()
		r.PrepTime = s.Find("[itemprop='prepTime']").Text()
		r.TotalTime = s.Find("[itemprop='totalTime']").Text()
		r.Yield = s.Find("[itemprop='recipeYield']").Text()
		r.Cuisine = s.Find("[itemprop='recipeCuisine']").Text()

		s.Find("[itemprop='recipeIngredient']").Each(func(i int, s *goquery.Selection) {
			r.Ingredients = append(r.Ingredients, strings.TrimSpace(s.Text()))
		})

		s.Find("[itemprop='recipeInstructions']").Each(func(i int, s *goquery.Selection) {
			r.Instructions = append(r.Instructions, strings.TrimSpace(s.Text()))
		})
	})

	if r.Title == "" {
		return nil, fmt.Errorf("no recipe found at URL")
	}

	return r, nil
}
