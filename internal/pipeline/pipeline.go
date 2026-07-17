package pipeline

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type ToolsConfig struct {
	FastQCPath         string
	UnicyclerPath      string
	SpadesPath         string
	CheckMPath         string
	Kraken2Path        string
	KrakenDBPath       string
	FastANIPath        string
	AbricatePath       string
	MLSTPath           string
	ResfinderDBPath    string
	PoliDbPseudo       string
	PoliDbKleb         string
	PoliDbEntero       string
	PoliDbAcineto      string
	OtherDbPseudo      string
	OtherDbKleb        string
	OtherDbEntero      string
	OtherDbAcineto     string
	FastaniListKleb    string
	FastaniListEntero  string
	FastaniListAcineto string
}

type CabgenPipeline interface {
	GetConfig() *ToolsConfig
	RunFastQC(ctx context.Context, read1, read2, outputDir string) (
		string, string, error)
	RunUnicycler(ctx context.Context, threads int,
		read1, read2, spadesPath, outputDir string) (string, error)
	RunProkka(ctx context.Context, threads int,
		assembly, outputDir string) error
	RunCheckM(ctx context.Context, threads int, sample, assemblyDir,
		outputDir string) (*CheckMResult, error)
	RunKraken2(ctx context.Context, threads int, assembly,
		outputDir string) (*KrakenSpecies, *KrakenSpecies, error)
	RunBlastX(ctx context.Context, query, DB, outputFile string) error
	RunAbricate(ctx context.Context, threads int, db, input,
		outputFile string) error
	ProcessSpecies(ctx context.Context, threads int,
		sampleID, mostCommon, assemblyPath, outputDir string) (
		*SpeciesResult, error)
}

type cabgenPipeline struct {
	Runner ToolRunner
	Config ToolsConfig
}

func NewCabgenPipeline(runner ToolRunner, config ToolsConfig) CabgenPipeline {
	return &cabgenPipeline{
		Runner: runner,
		Config: config,
	}
}

func (p *cabgenPipeline) GetConfig() *ToolsConfig {
	return &p.Config
}

func (p *cabgenPipeline) RunFastQC(
	ctx context.Context, read1, read2, outputDir string) (
	string, string, error) {
	fastqcCmdArgs := p.Runner.BuildFastQCCmd(p.Config.FastQCPath, read1, read2,
		outputDir)

	if _, err := p.Runner.Run(ctx, fastqcCmdArgs); err != nil {
		return "", "", err
	}

	outputHTMLfile1 := filepath.Join(outputDir,
		fmt.Sprintf("%s_fastqc.html", filepath.Base(read1)))
	outputHTMLfile2 := filepath.Join(outputDir,
		fmt.Sprintf("%s_fastqc.html", filepath.Base(read2)))

	return outputHTMLfile1, outputHTMLfile2, nil
}

func (p *cabgenPipeline) RunUnicycler(ctx context.Context, threads int,
	read1, read2, spadesPath, outputDir string) (string, error) {
	threadsStr := strconv.Itoa(threads)

	unicyclerCmdArgs := p.Runner.BuildUnicyclerCmd(
		p.Config.UnicyclerPath, read1, read2, outputDir, threadsStr,
		p.Config.SpadesPath)

	if _, err := p.Runner.Run(ctx, unicyclerCmdArgs); err != nil {
		return "", err
	}

	assemblyPath := filepath.Join(outputDir, "assembly.fasta")

	return assemblyPath, nil
}

func (p *cabgenPipeline) RunProkka(ctx context.Context, threads int,
	assembly, outputDir string) error {
	threadsStr := strconv.Itoa(threads)

	prokkaCmd := "prokka"
	prefix := "genome"
	prokkaCmdArgs := p.Runner.BuildProkkaCmd(prokkaCmd, outputDir,
		prefix, assembly, threadsStr)

	if _, err := p.Runner.Run(ctx, prokkaCmdArgs); err != nil {
		return err
	}

	return nil
}

func (p *cabgenPipeline) RunCheckM(ctx context.Context, threads int,
	sample, assemblyDir, outputDir string) (*CheckMResult, error) {
	threadsStr := strconv.Itoa(threads)

	lineageArgs := p.Runner.BuildCheckMLineageCmd(p.Config.CheckMPath,
		assemblyDir, outputDir, threadsStr)
	if _, err := p.Runner.Run(ctx, lineageArgs); err != nil {
		return nil, err
	}

	qaArgs := p.Runner.BuildCheckMQACmd(p.Config.CheckMPath, outputDir,
		sample, threadsStr)
	if _, err := p.Runner.Run(ctx, qaArgs); err != nil {
		return nil, err
	}

	resultPath := filepath.Join(outputDir, fmt.Sprintf("%s_results", sample))
	result, err := ParseCheckM(resultPath)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *cabgenPipeline) RunKraken2(ctx context.Context, threads int, assembly,
	outputDir string) (*KrakenSpecies, *KrakenSpecies, error) {
	threadsStr := strconv.Itoa(threads)

	krakenArgs := p.Runner.BuildKraken2Cmd(
		p.Config.Kraken2Path, p.Config.KrakenDBPath, outputDir,
		threadsStr, assembly,
	)
	if _, err := p.Runner.Run(ctx, krakenArgs); err != nil {
		return nil, nil, err
	}

	krakenResult := filepath.Join(outputDir, "out_kraken")
	kResult1, kResult2, err := KrakenSpeciesCounter(krakenResult)
	if err != nil {
		return nil, nil, err
	}

	return kResult1, kResult2, nil
}

func (p *cabgenPipeline) RunBlastX(ctx context.Context, query, DB,
	outputFile string) error {
	blastArgs := p.Runner.BuildBlastXCmd(DB, query, outputFile)
	if _, err := p.Runner.Run(ctx, blastArgs); err != nil {
		return err
	}

	return nil
}

func (p *cabgenPipeline) RunAbricate(ctx context.Context, threads int, db,
	input, outputFile string) error {
	threadsStr := strconv.Itoa(threads)

	abricateArgs := p.Runner.BuildAbricateCmd(p.Config.AbricatePath, db,
		input, outputFile, threadsStr)
	if _, err := p.Runner.Run(ctx, abricateArgs); err != nil {
		return err
	}

	return nil
}

func (p *cabgenPipeline) ProcessSpecies(ctx context.Context, threads int,
	sampleID, mostCommon, assemblyPath, outputDir string) (
	*SpeciesResult, error) {
	mostCommon = strings.TrimSpace(mostCommon)
	parts := strings.Fields(mostCommon)

	genus := mostCommon
	species := ""
	if len(parts) >= 2 {
		genus = parts[0]
		species = parts[1]
	}

	normalizedName := strings.ToLower(genus + species)

	displayName := fmt.Sprintf("%s %s", capitalizeFirst(genus),
		strings.ToLower(species))
	if species == "" {
		displayName = capitalizeFirst(genus)
	}

	result := &SpeciesResult{
		DisplayName:    strings.TrimSpace(displayName),
		MLSTSpecies:    "",
		OtherMutations: []string{},
		PoliMutations:  []string{},
	}

	threadsStr := strconv.Itoa(threads)

	mlstResultPath := filepath.Join(outputDir, "mlst.csv")
	mlstArgs := p.Runner.BuildMLSTCmd(p.Config.MLSTPath, threadsStr,
		assemblyPath, mlstResultPath)
	if _, err := p.Runner.Run(ctx, mlstArgs); err == nil {
		if mlstData, err := ParseMLST(mlstResultPath); err == nil &&
			mlstData != nil {
			result.MLSTSpecies = fmt.Sprintf(
				"%s (ST: %s)", mlstData.Scheme, mlstData.ST)
		}
	}

	isEntero := isEnterobacter(normalizedName)
	isAcineto := isAcinetobacter(normalizedName)
	isKleb := isKlebsiella(normalizedName)
	isPseudo := isPseudomonas(normalizedName)

	var poliDbFullPath, otherDbFullPath, fastAniRefFullPath string

	if isPseudo {
		poliDbFullPath = p.Config.PoliDbPseudo
		otherDbFullPath = p.Config.OtherDbPseudo
	} else if isKleb {
		poliDbFullPath = p.Config.PoliDbKleb
		otherDbFullPath = p.Config.OtherDbKleb
		fastAniRefFullPath = p.Config.FastaniListKleb
	} else if isEntero {
		poliDbFullPath = p.Config.PoliDbEntero
		otherDbFullPath = p.Config.OtherDbEntero
		fastAniRefFullPath = p.Config.FastaniListEntero
	} else if isAcineto {
		poliDbFullPath = p.Config.PoliDbAcineto
		otherDbFullPath = p.Config.OtherDbAcineto
		fastAniRefFullPath = p.Config.FastaniListAcineto
	}

	if (isEntero || isAcineto || isKleb) && fastAniRefFullPath != "" {
		fastAniOut := filepath.Join(outputDir,
			fmt.Sprintf("%s_out-fastANI", sampleID))

		fastAniArgs := p.Runner.BuildFastANICmd(
			p.Config.FastANIPath, assemblyPath, fastAniRefFullPath,
			fastAniOut, threadsStr,
		)
		if _, err := p.Runner.Run(ctx, fastAniArgs); err == nil {
			fastAniSpecies, parseErr := ParseFastANI(fastAniOut)
			if parseErr == nil && fastAniSpecies != "" {
				result.DisplayName = strings.ReplaceAll(fastAniSpecies, "_",
					" ")
			}
		}
	}

	if poliDbFullPath != "" && otherDbFullPath != "" {
		blastPoliFile := filepath.Join(outputDir, fmt.Sprintf(
			"%s_blastPoli", sampleID))

		if err := p.RunBlastX(ctx, assemblyPath, poliDbFullPath,
			blastPoliFile); err != nil {
			return nil, err
		}

		blastOtherFile := filepath.Join(outputDir, fmt.Sprintf(
			"%s_blastOther", sampleID))

		if err := p.RunBlastX(ctx, assemblyPath, otherDbFullPath,
			blastOtherFile); err != nil {
			return nil, err
		}

		poliFinder := NewMutationFinder(blastPoliFile)
		otherFinder := NewMutationFinder(blastOtherFile)

		var otherMut, poliMut []string
		var errPoli, errOther error

		if isAcineto {
			_, poliMut, errPoli = poliFinder.FindAcinetoMutations()
			otherMut, _, errOther = otherFinder.FindAcinetoMutations()
		} else if isEntero {
			_, poliMut, errPoli = poliFinder.FindEcloacaeMutations()
			otherMut, _, errOther = otherFinder.FindEcloacaeMutations()
		} else if isKleb {
			_, poliMut, errPoli = poliFinder.FindKlebMutations()
			otherMut, _, errOther = otherFinder.FindKlebMutations()
		} else if isPseudo {
			_, poliMut, errPoli = poliFinder.FindPseudoMutations()
			otherMut, _, errOther = otherFinder.FindPseudoMutations()
		}

		if errPoli == nil && poliMut != nil {
			result.PoliMutations = poliMut
		}
		if errOther == nil && otherMut != nil {
			result.OtherMutations = otherMut
		}
	}

	return result, nil
}
