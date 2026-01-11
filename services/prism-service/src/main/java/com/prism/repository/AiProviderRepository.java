package com.prism.repository;

import com.prism.domain.AiProvider;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface AiProviderRepository extends JpaRepository<AiProvider, String> {
}
